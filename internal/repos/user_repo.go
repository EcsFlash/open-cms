package repos

import (
	"headless-cms/internal/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

type IUserRepo interface {
	Create(u *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
}

func (r *UserRepo) Create(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepo) GetByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByNickname(nickname string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("nickname = ?", nickname).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

