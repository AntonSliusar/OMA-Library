package server

import (
	"net/http"
	"oma-library/internal/config"
	"oma-library/internal/handlers"
	"oma-library/internal/utils"

	"github.com/go-playground/validator/v10"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RunServer(cfg *config.Config, fileHandler *handlers.FileHandler, userHandler *handlers.UserHandler) { //change arguments to FileHandler
	e := echo.New()
	e.Renderer = utils.NewTemplate("templates/*.html")
	e.Static("/", "templates")
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	resGroup := e.Group("/admin")
	resGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  cfg.AUTH.JWT,
		TokenLookup: "cookie:jwt",
	}))

	e.GET("/", handlers.MainPageHandler)
	e.GET("/search", fileHandler.SearchHandler)
	e.GET("/oma/:id", fileHandler.DownloadHandler)
	e.POST("/register", userHandler.RegisterUser)
	e.POST("/login", userHandler.UserLogin)
	resGroup.GET("/upload", handlers.UploadFormHandler)
	resGroup.POST("/upload_file", fileHandler.UploadFileHandler)
	resGroup.GET("/check", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "valid"})
	})

	e.Start(":8080") // add to config

}