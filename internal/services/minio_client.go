package services

import (
	"context"
	"headless-cms/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIO(cfg *config.Config) (*minio.Client, error) {
	return minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
}

func EnsureBucket(ctx context.Context, cli *minio.Client, bucket string) error {
	ok, err := cli.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return cli.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
}

