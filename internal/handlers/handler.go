package handlers

import (
	"errors"
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
	e.GET("/search", searchHandler(storage, r2))
	e.GET("/oma/:id", dowloadHandler(storage, r2))
	e.POST("/register", RegisterAdmin(storage))
	e.POST("/login", AdminLogin(storage))
	resGroup.GET("/upload", UploadFormHandler)
	resGroup.POST("/upload_file", UploadFileHandler(storage, r2))
	resGroup.GET("/check", func(c echo.Context) error {
    	return c.JSON(http.StatusOK, map[string]string{"status": "valid"})
	})

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
		imgHeader, err := c.FormFile("image")
		if err != nil {
			logger.Logger.Info(errors.New("no image to upload"))
		}


		if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".oma") {
			logger.Logger.Info("Not oma")
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid file format. Only .oma files are allowed"})

		}
		omaFile := new(models.Omafile)
		omaFile.Brand = strings.ToLower(c.FormValue("Brand"))
		omaFile.Model = strings.ToLower(c.FormValue("Model"))
		omaFile.OMAKey = utils.AddPrefix(fileHeader.Filename)
		if imgHeader != nil {
			omaFile.ImgKey = utils.AddPrefix(imgHeader.Filename)
		}

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
		err = r2.UploadFileToR2(c.Request().Context(), omaFile.OMAKey, file)
		if err != nil {
			logger.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Cannot upload file to R2"})
		}

		if imgHeader != nil {
			imgFile, err := imgHeader.Open()
			if err != nil {
				logger.Logger.Error(err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Cannot open image"})
			}
			err = r2.UploadFileToR2(c.Request().Context(), omaFile.ImgKey, imgFile)
			if err != nil {
				logger.Logger.Error(err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Cannot upload image to R2"})
			}
			defer imgFile.Close()
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "File uploaded successfully"})
	}
}

func searchHandler(storage *storage.Storage, r2 *storage.R2Client) echo.HandlerFunc {
	return func(c echo.Context) error {

		brand := strings.ToLower(c.QueryParam("brand"))
		model := strings.ToLower(c.QueryParam("model"))
		var files []models.Omafile = nil

		if brand != "" && model != "" {
			files = storage.GetByBrandAndModel(brand, model)
		} else if brand != "" {
			files = storage.GetByBrand(brand)
		} else if model != "" {
			files = storage.GetByModel(model)
		}

		for i := range files {
			if files[i].ImgKey != "" {
				url, err := r2.GeneratePresignedURLForImg(c.Request().Context(), files[i].ImgKey)
				if err != nil {
					logger.Logger.Error(err)
					continue
				}
				files[i].ImgURL = url
			}
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
		obj, err := r2.DownloadFileFromR2(c.Request().Context(), oma.OMAKey)
		if err != nil {
			logger.Logger.Error(err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		defer obj.Body.Close()


		c.Response().Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(oma.OMAKey))
		c.Response().Header().Set("Content-Type", "application/octet-stream")

		_, err = io.Copy(c.Response().Writer, obj.Body)
		if err != nil {
			logger.Logger.Error(err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		return nil
	}
}
