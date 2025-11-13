package entity

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int
	Username string
	PasswordHash string
}

func (u * User) CheckPassword(raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(raw)) == nil
}