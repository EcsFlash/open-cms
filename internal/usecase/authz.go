package usecase

import (
	"errors"
	"headless-cms/internal/models"
)

var ErrForbidden = errors.New("forbidden")

type Actor struct {
	UserID uint
	Role   models.Role
}

func (a Actor) IsAdmin() bool {
	return a.Role == models.RoleAdmin
}

func EnsureCanMutate(actor Actor, ownerID uint) error {
	if actor.IsAdmin() || actor.UserID == ownerID {
		return nil
	}
	return ErrForbidden
}
