package handlers

import (
	"io"
	"net/http"
	"oma-library/internal/models"
	"oma-library/internal/utils"
	"oma-library/pkg/logger"
	"oma-library/pkg/storage"
	"path/filepath"

	"strings"

	"github.com/go-playground/validator/v10"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func RunWeb(storage *storage.Storage, r2 *storage.R2Client) {
	e := echo.New()
	e.Renderer = utils.NewTemplate("templates/*.html")
	e.Static("/", "templates")
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	resGroup := e.Group("/admin")
	resGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  jwtSecret,
		TokenLookup: "cookie:jwt",
	}))

	e.GET("/", MainPageHandler)
	e.GET("/search", searchHandler(storage))
	e.GET("/oma/:id", dowloadHandler(storage, r2))
	e.POST("/register", RegisterAdmin(storage))
	e.POST("/login", AdminLogin(storage))
	resGroup.GET("/upload", UploadFormHandler)
	resGroup.POST("/upload_file", UploadFileHandler(storage, r2))

	e.Start(":8080") // add to config

}

func MainPageHandler(c echo.Context) error {
	if err := c.Render(http.StatusOK, "index", nil); err != nil {
		logger.Logger.Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	return nil
}

func UploadFormHandler(c echo.Context) error {
	if err := c.Render(http.StatusOK, "upload", nil); err != nil {
		logger.Logger.Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	return nil
}

// hendlers working with db:

func UploadFileHandler(storage *storage.Storage, r2 *storage.R2Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		fileHeader, err := c.FormFile("uploaded_file")
		if err != nil {
			logger.Logger.Error(err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "File is required"})
		}


		if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".oma") {
			logger.Logger.Info("Not oma")
			c.String(http.StatusBadRequest, "Invalid file format. Only .oma files are allowed"+err.Error())
		}
		omaFile := new(models.Omafile)
		omaFile.Brand = c.FormValue("Brand")
		omaFile.Model = c.FormValue("Model")
		omaFile.Key = utils.AddPrefix(fileHeader.Filename)

		err = storage.Create(*omaFile)
		if err != nil {
			logger.Logger.Error(err)
		}

		file, err := fileHeader.Open()
		if err != nil {
			logger.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Cannot open uploaded file"})
		}
		defer file.Close()
		
		err = r2.UploadFileToR2(c.Request().Context(), omaFile.Key, file)
		if err != nil {
			logger.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Cannot upload file to R2"})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "File uploaded successfully"})
	}
}

func searchHandler(storage *storage.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {

		brand := c.QueryParam("brand")
		model := c.QueryParam("model")
		var files []models.Omafile = nil

		if brand != "" && model != "" {
			files = storage.GetByBrandAndModel(brand, model)
		} else if brand != "" {
			files = storage.GetByBrand(brand)
		} else if model != "" {
			files = storage.GetByModel(model)
		}

		if err := c.Render(http.StatusOK, "forms", files); err != nil {
			logger.Logger.Error(err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		return nil
	}
}

func dowloadHandler(storage *storage.Storage, r2 *storage.R2Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		logger.Logger.Println(id)
		oma := storage.GetById(id)
		obj, err := r2.DownloadFileFromR2(c.Request().Context(), oma.Key)
		if err != nil {
			logger.Logger.Error(err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		defer obj.Body.Close()


		c.Response().Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(oma.Key))
		c.Response().Header().Set("Content-Type", "application/octet-stream")

		_, err = io.Copy(c.Response().Writer, obj.Body)
		if err != nil {
			logger.Logger.Error(err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		return nil
	}
}
