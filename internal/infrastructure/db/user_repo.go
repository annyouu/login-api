package db

import (
	"database/sql"
	"login/internal/domain/entity"
	repo "login/internal/repository"
)

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) repo.UserRepositorInterface {
	return &userRepository{DB:db}
}

func (r *userRepository) GetByID(id int64) (*entity.User, error) {
	user := &entity.User{}
	query := "SELECT id, username, password_hash FROM users WHERE id = ? LIMIT 1"

	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByUsername(username string) (*entity.User, error) {
	user := &entity.User{}
	query := "SELECT id, username, password_hash FROM users WHERE username = ? LIMIT 1"

	err := r.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Create(user *entity.User) error {
    query := "INSERT INTO users (username, password_hash) VALUES (?, ?)"
    result, err := r.DB.Exec(query, user.Username, user.PasswordHash)
    if err != nil {
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return err
    }

    user.ID = id
    return nil
}

