package repository

import "login/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	GetByUsername(username string) (*entity.User, error)
	GetByID(id int64) (*entity.User, error)
}