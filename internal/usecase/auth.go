package usecase

import (
	"errors"
	"login/internal/domain/entity"
	"login/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecaseInterface interface {
	Register(username, password string) error
    Login(username, password string) (string, error)
    GetUser(id int64) (*entity.User, error)
}

type authUsecase struct {
	userRepo repository.UserRepositorInterface 
	jwtGenerator JWTGeneratorInterface
}

func NewAuthUsecase(
	userRepo repository.UserRepositorInterface ,
	jwtGen JWTGeneratorInterface,
) AuthUsecaseInterface {
	return &authUsecase{
		userRepo: userRepo,
		jwtGenerator: jwtGen,
	}
}

func (u *authUsecase) Register(username, password string) error {
	// 既存ユーザーを確認する
	existing, _ := u.userRepo.GetByUsername(username)
	if existing != nil {
		return errors.New("ユーザーはすでに存在しています")
	}

	// パスワードをハッシュ化する
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.User{
		Username: username,
		PasswordHash: string(hash),
	}

	// 永続化 (DBに保存する)
	return u.userRepo.Create(user)
}

func (u *authUsecase) Login(username, password string) (string, error) {
	user, err := u.userRepo.GetByUsername(username)
	if err != nil || user == nil {
		return "", errors.New("ユーザーが存在しません")
	}

	// パスワード一致チェックする
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("パスワードが間違っています")
	}

	// JWT生成
	token, err := u.jwtGenerator.Generate(user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *authUsecase) GetUser(id int64) (*entity.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("ユーザーが見つかりません")
	}
	return user, nil
}
