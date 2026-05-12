package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"headless-cms/internal/config"
	"headless-cms/internal/models"
	"headless-cms/internal/repos"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNicknameTaken      = errors.New("nickname already taken")
)

type AuthService struct {
	cfg  *config.Config
	repo repos.IUserRepo
}

type IAuthService interface {
	Register(nickname, password string) (*models.User, error)
	Login(nickname, password string) (string, *models.User, error)
	RegisterSupervisor(nickname, password string, role models.Role) (*models.User, error)
	RemoveSupervisor(nickname string) error
}

func NewAuthService(cfg *config.Config, repo repos.IUserRepo) *AuthService {
	if _, err := repo.GetByNickname(cfg.Admin.Nickname); err != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Admin.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		u := &models.User{
			Nickname:     cfg.Admin.Nickname,
			PasswordHash: string(hash),
			Role:         models.RoleAdmin,
		}
		if err := repo.Create(u); err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
			panic(err)
		}
	}
	return &AuthService{cfg: cfg, repo: repo}
}

func (s *AuthService) Register(nickname, password string) (*models.User, error) {
	if nickname == "" || password == "" {
		return nil, errors.New("nickname and password are required")
	}
	if _, err := s.repo.GetByNickname(nickname); err == nil {
		return nil, ErrNicknameTaken
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &models.User{
		Nickname:     nickname,
		PasswordHash: string(hash),
		Role:         models.RoleUser,
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthService) RegisterSupervisor(nickname, password string, role models.Role) (*models.User, error) {
	if nickname == "" || password == "" {
		return nil, errors.New("nickname and password are required")
	}
	if _, err := s.repo.GetByNickname(nickname); err == nil {
		return nil, ErrNicknameTaken
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &models.User{
		Nickname:     nickname,
		PasswordHash: string(hash),
		Role:         role,
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthService) Login(nickname, password string) (string, *models.User, error) {
	u, err := s.repo.GetByNickname(nickname)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  u.ID,
		"role": string(u.Role),
		"iat":  now.Unix(),
		"exp":  now.Add(72 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", nil, err
	}
	return token, u, nil
}

func (s *AuthService) RemoveSupervisor(nickname string) error {
	u, err := s.repo.GetByNickname(nickname)
	if err != nil {
		return err
	}
	if u.Role != models.RoleModerator {
		return fmt.Errorf("%s is not a moderator", nickname)
	}
	err = s.repo.Remove(u)
	if err != nil {
		return err
	}
	return nil
}
