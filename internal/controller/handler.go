package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"login/internal/usecase"
)

type UserHandler struct {
	authUC usecase.AuthUsecaseInterface
}

func NewUserHandler(authUC usecase.AuthUsecaseInterface) *UserHandler {
	return &UserHandler{
		authUC: authUC,
	}
}

// POST/register
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "無効なリクエストです"})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "名前、パスワードは必須です"})
	}

	err := h.authUC.Register(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "ユーザー登録が完了しました"})
}

// POST/login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	token, err := h.authUC.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}

	// CookieにJWTをセットする
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{"message": "ログインが成功しました"})
}

// GET /users/id
func (h *UserHandler) GetUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	user, err := h.authUC.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id": user.ID,
		"username": user.Username,
	})
}