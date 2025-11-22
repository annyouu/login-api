package usecase_test

import (
	"errors"
	"login/internal/domain/entity"
	"login/internal/usecase"
	"testing"
)

// ---- モック定義 ----

type mockUserRepo struct {
	users map[string]*entity.User
}

func (m *mockUserRepo) Create(user *entity.User) error {
	m.users[user.Username] = user
	return nil
}

func (m *mockUserRepo) GetByUsername(username string) (*entity.User, error) {
	if u, ok := m.users[username]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) GetByID(id int64) (*entity.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

// JWT のモック
type mockJWT struct{}

func (m *mockJWT) Generate(username string) (string, error) {
	return "mock_token_for_" + username, nil
}

// ---- テスト開始 ----

func TestRegisterSuccess(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]*entity.User)}
	jwt := &mockJWT{}
	uc := usecase.NewAuthUsecase(repo, jwt)

	err := uc.Register("testuser", "password")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.users["testuser"] == nil {
		t.Fatalf("user should be created")
	}
}

func TestRegisterDuplicateUser(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]*entity.User)}
	repo.users["testuser"] = &entity.User{Username: "testuser"}
	jwt := &mockJWT{}
	uc := usecase.NewAuthUsecase(repo, jwt)

	err := uc.Register("testuser", "password")

	if err == nil {
		t.Fatalf("expected error for duplicate user")
	}
}

func TestLoginSuccess(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]*entity.User)}
	repo.users["testuser"] = &entity.User{
		Username:     "testuser",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZo5e.Kl8QJHki71R.plVQ1BMdDT6vKrrO5e.", // bcrypt("password")
	}
	jwt := &mockJWT{}
	uc := usecase.NewAuthUsecase(repo, jwt)

	token, err := uc.Login("testuser", "password")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token != "mock_token_for_testuser" {
		t.Fatalf("unexpected token: %v", token)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	repo := &mockUserRepo{users: make(map[string]*entity.User)}
	repo.users["testuser"] = &entity.User{
		Username:     "testuser",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZo5e.Kl8QJHki71R.plVQ1BMdDT6vKrrO5e.", // password
	}
	jwt := &mockJWT{}
	uc := usecase.NewAuthUsecase(repo, jwt)

	_, err := uc.Login("testuser", "wrong")

	if err == nil {
		t.Fatalf("expected error for wrong password")
	}
}
