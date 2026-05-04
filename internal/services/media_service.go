package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"headless-cms/internal/config"
	"headless-cms/internal/models"
	"headless-cms/internal/repos"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type MediaService struct {
	cfg    *config.Config
	minio  *minio.Client
	images repos.IImageRepo
	videos repos.IVideoRepo
}

type IMediaService interface {
	UploadImage(ctx context.Context, file multipart.File, filename, contentType, name string) (*models.Image, error)
	UploadVideo(ctx context.Context, file multipart.File, filename, contentType, name string, thumbnailImageID *uint) (*models.Video, error)

	GetImageByID(id uint) (*models.Image, error)
	GetVideoByID(id uint) (*models.Video, error)
	ListImages() ([]models.Image, error)
	ListVideos() ([]models.Video, error)

	RenameImage(id uint, name string) error
	DeleteImage(ctx context.Context, id uint) error
	RenameVideo(id uint, name string) error
	DeleteVideo(ctx context.Context, id uint) error
}

func NewMediaService(cfg *config.Config, minioCli *minio.Client, images repos.IImageRepo, videos repos.IVideoRepo) *MediaService {
	return &MediaService{cfg: cfg, minio: minioCli, images: images, videos: videos}
}

func (s *MediaService) UploadImage(ctx context.Context, file multipart.File, filename, contentType, name string) (*models.Image, error) {
	if name == "" {
		name = filename
	}
	data, err := readAll(file, 25<<20) // 25MB soft limit for now
	if err != nil {
		return nil, err
	}
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}
	_ = format
	b := img.Bounds()

	base := uuid.NewString()
	origKey := fmt.Sprintf("images/%s/original%s", base, filepath.Ext(filename))
	if strings.TrimSpace(filepath.Ext(filename)) == "" {
		origKey = fmt.Sprintf("images/%s/original", base)
	}

	_, err = s.minio.PutObject(ctx, s.cfg.MinIO.ImagesBucket, origKey, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	variant := func(key string, width int, quality int) (string, error) {
		resized := imaging.Resize(img, width, 0, imaging.Lanczos)
		var buf bytes.Buffer
		if err := imaging.Encode(&buf, resized, imaging.JPEG, imaging.JPEGQuality(quality)); err != nil {
			return "", err
		}
		vKey := fmt.Sprintf("images/%s/%s.jpg", base, key)
		_, err := s.minio.PutObject(ctx, s.cfg.MinIO.ImagesBucket, vKey, bytes.NewReader(buf.Bytes()), int64(buf.Len()), minio.PutObjectOptions{
			ContentType: "image/jpeg",
		})
		if err != nil {
			return "", err
		}
		return vKey, nil
	}

	largeKey, err := variant("large", 1600, 85)
	if err != nil {
		return nil, err
	}
	mediumKey, err := variant("medium", 800, 80)
	if err != nil {
		return nil, err
	}
	miniKey, err := variant("mini", 400, 75)
	if err != nil {
		return nil, err
	}
	thumbKey, err := variant("thumbnail", 200, 70)
	if err != nil {
		return nil, err
	}

	m := &models.Image{
		Name:               name,
		Mime:               contentType,
		Width:              b.Dx(),
		Height:             b.Dy(),
		Bucket:             s.cfg.MinIO.ImagesBucket,
		ObjectKeyOriginal:  origKey,
		ObjectKeyLarge:     largeKey,
		ObjectKeyMedium:    mediumKey,
		ObjectKeyMini:      miniKey,
		ObjectKeyThumbnail: thumbKey,
	}
	if err := s.images.Create(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MediaService) UploadVideo(ctx context.Context, file multipart.File, filename, contentType, name string, thumbnailImageID *uint) (*models.Video, error) {
	if name == "" {
		name = filename
	}
	tmp, err := os.CreateTemp("", "cms-video-*"+filepath.Ext(filename))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
	}()

	size, err := copyToFile(tmp, file, 2<<30) // 2GB limit
	if err != nil {
		return nil, err
	}

	if contentType == "" {
		// best-effort: rely on client for videos; otherwise generic
		contentType = "application/octet-stream"
	}

	base := uuid.NewString()
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".mp4"
	}
	videoKey := fmt.Sprintf("videos/%s/original%s", base, ext)
	_, err = s.minio.FPutObject(ctx, s.cfg.MinIO.VideosBucket, videoKey, tmp.Name(), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	duration, _ := probeDurationSeconds(tmp.Name())

	thumbID := thumbnailImageID
	if thumbID == nil {
		thumbPath := tmp.Name() + ".thumb.jpg"
		if err := generateVideoThumbnail(tmp.Name(), thumbPath); err == nil {
			defer os.Remove(thumbPath)

			f, err := os.Open(thumbPath)
			if err == nil {
				defer f.Close()
				imgModel, err := s.UploadImage(ctx, f, "thumbnail.jpg", "image/jpeg", name+" thumbnail")
				if err == nil {
					thumbID = &imgModel.ID
				}
			}
		}
	}

	v := &models.Video{
		Name:              name,
		Mime:              contentType,
		Bucket:            s.cfg.MinIO.VideosBucket,
		ObjectKey:         videoKey,
		Duration:          duration,
		ThumbnailImageID:  thumbID,
		ThumbnailImage:    nil,
	}
	if err := s.videos.Create(v); err != nil {
		return nil, err
	}
	_ = size
	return v, nil
}

func (s *MediaService) GetImageByID(id uint) (*models.Image, error) { return s.images.GetByID(id) }
func (s *MediaService) GetVideoByID(id uint) (*models.Video, error) { return s.videos.GetByID(id) }
func (s *MediaService) ListImages() ([]models.Image, error)         { return s.images.ListAll() }
func (s *MediaService) ListVideos() ([]models.Video, error)         { return s.videos.ListAll() }

func (s *MediaService) RenameImage(id uint, name string) error {
	return s.images.UpdateName(id, name)
}

func (s *MediaService) DeleteImage(ctx context.Context, id uint) error {
	img, err := s.images.GetByID(id)
	if err != nil {
		return err
	}
	for _, key := range []string{
		img.ObjectKeyOriginal,
		img.ObjectKeyLarge,
		img.ObjectKeyMedium,
		img.ObjectKeyMini,
		img.ObjectKeyThumbnail,
	} {
		_ = s.minio.RemoveObject(ctx, img.Bucket, key, minio.RemoveObjectOptions{})
	}
	return s.images.Delete(id)
}

func (s *MediaService) RenameVideo(id uint, name string) error {
	return s.videos.UpdateName(id, name)
}

func (s *MediaService) DeleteVideo(ctx context.Context, id uint) error {
	v, err := s.videos.GetByID(id)
	if err != nil {
		return err
	}
	_ = s.minio.RemoveObject(ctx, v.Bucket, v.ObjectKey, minio.RemoveObjectOptions{})
	return s.videos.Delete(id)
}

func readAll(r multipart.File, max int64) ([]byte, error) {
	lr := io.LimitReader(r, max+1)
	b, err := io.ReadAll(lr)
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > max {
		return nil, errors.New("file too large")
	}
	return b, nil
}

func copyToFile(dst *os.File, src multipart.File, max int64) (int64, error) {
	lr := io.LimitReader(src, max+1)
	n, err := io.Copy(dst, lr)
	if err != nil {
		return n, err
	}
	if n > max {
		return n, errors.New("file too large")
	}
	if _, err := dst.Seek(0, 0); err != nil {
		return n, err
	}
	return n, nil
}

func probeDurationSeconds(path string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", path)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return 0, errors.New("empty duration")
	}
	return strconv.ParseFloat(s, 64)
}

func generateVideoThumbnail(videoPath, outputPath string) error {
	// First frame (t=0) scaled to 1280 width.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ffmpeg", "-y", "-i", videoPath, "-ss", "0", "-vframes", "1", "-q:v", "18", "-vf", "scale=1280:-1", outputPath)
	return cmd.Run()
}

