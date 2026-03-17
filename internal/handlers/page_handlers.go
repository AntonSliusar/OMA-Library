package handlers

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func MainPageHandler(c echo.Context) error {
	if err := c.Render(http.StatusOK, "index", nil); err != nil {
		slog.Error(err.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	return nil
}

func UploadFormHandler(c echo.Context) error {
	if err := c.Render(http.StatusOK, "upload", nil); err != nil {
		slog.Error(err.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	return nil
}