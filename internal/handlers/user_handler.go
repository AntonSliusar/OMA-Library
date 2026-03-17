package handlers

import (
	"net/http"
	"oma-library/internal/models"
	"oma-library/internal/service"
	"oma-library/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("secret_token") //Must be at config??


type UserHandler struct {
	service *service.UserService
}

func NewUserHandler (s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func(h *UserHandler) RegisterUser(c echo.Context) error {
	var signUpReq models.SignUpRequest
	err := c.Bind(&signUpReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "wrong data"})
	}
	err = c.Validate(signUpReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}
	exist, err := h.service.CheckExist(signUpReq.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{})
	}
	if exist {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Email already exist"})
	}
	signUpReq.Password = utils.HashPassword(signUpReq.Password)
	err = h.service.RegisterUser(signUpReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "cant create admin"})
	}
	return c.JSON(http.StatusCreated, nil)
}

func(h UserHandler) UserLogin(c echo.Context) error {
	var signInReq models.SignInRequset
	err := c.Bind(&signInReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "wrong data"})
	}
	err = c.Validate(signInReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	user, err := h.service.LoginUser(signInReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	if utils.CheckPasswordHash(signInReq.Password, user.Password) != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Could not generate token"})
	}
	c.SetCookie(&http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		MaxAge: 86400,
	})
	return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
}
