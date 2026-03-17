package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"oma-library/internal/service"
	"path/filepath"

	"strings"

	"github.com/labstack/echo/v4"
)

type FileHandler struct {
	service *service.OmaFileService
}

func NewFileHandler(service *service.OmaFileService) *FileHandler {
	return &FileHandler{service: service}
}



// hendlers working with db:

func (h *FileHandler) UploadFileHandler(c echo.Context) error {
	var input service.UploadFileInput
	fileHeader, err := c.FormFile("uploaded_file")
	if err != nil {
		slog.Error(err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "File is required"})
	}
	imgHeader, _ := c.FormFile("image")
	input.File = fileHeader
	input.Image = imgHeader
	input.Brand = c.FormValue("Brand")
	input.Model = c.FormValue("Model")

	err = h.service.UploadOmaFile(c.Request().Context(), input)
	if err != nil {
		if err == service.ErrInvalidFileFormat {
			slog.Info("Not oma")
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid file format. Only .oma files are allowed"})
		}
		slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "service error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "File uploaded successfully"})		
}


func (h *FileHandler) SearchHandler(c echo.Context) error {
	brand := strings.ToLower(c.QueryParam("brand"))
	model := strings.ToLower(c.QueryParam("model"))
	files := h.service.SearchOmaFile(c.Request().Context(), brand, model)

	if err := c.Render(http.StatusOK, "forms", files); err != nil {
		slog.Error(err.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	return nil
}

func (h *FileHandler) DownloadHandler(c echo.Context) error {
	id := c.Param("id")
	output, err := h.service.DownloadFile(c.Request().Context(), id)
	if err != nil {
		slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "service error"})
	}
	defer func() {
		if output.File != nil && output.File.Body != nil {
			output.File.Body.Close()
		}
	}()
	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(output.Oma.OMAKey))
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(c.Response().Writer, output.File.Body)
	if err != nil {
		slog.Error(err.Error())
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	return nil
}


