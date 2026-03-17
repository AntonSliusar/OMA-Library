package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"oma-library/internal/models"
	"oma-library/internal/utils"
	"oma-library/pkg/storage"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// type FileStorage interface {
// 	UploadFile(ctx context.Context, file models.Omafile) error
// }

type OmaFileService struct {
	storage *storage.Storage
	r2 *storage.R2Client
}

type UploadFileInput struct {
	File     *multipart.FileHeader
    Image    *multipart.FileHeader 
    Brand    string
    Model    string
}

type DownloadFileOutput struct {
	Oma models.Omafile
	File *s3.GetObjectOutput
}

func NewOmaFileService(storage *storage.Storage, r2 *storage.R2Client) *OmaFileService {
	return &OmaFileService{
		storage: storage,
		r2: r2,
	}
}

func (s *OmaFileService) UploadOmaFile(ctx context.Context, input UploadFileInput) error {
	if !strings.HasSuffix(strings.ToLower(input.File.Filename), ".oma") {
			slog.Info("Not oma")
			return ErrInvalidFileFormat
		}

	omaFile := models.Omafile{}
	omaFile.Brand = strings.ToLower(input.Brand)
	omaFile.Model = strings.ToLower(input.Model)
	omaFile.OMAKey = utils.AddPrefix(input.File.Filename)
	if input.Image != nil {
		omaFile.ImgKey = utils.AddPrefix(input.Image.Filename)
	}

	id, err := s.storage.Create(omaFile)
	if err != nil {
		return fmt.Errorf("create file err: %w", err)
	}

	//upload to r2
	err = s.uploadToObjectStorage(ctx, input, omaFile)
	if err != nil {
		if errors.Is(err, ErrImgUpload) {
			return ErrImgUpload
		}
		s.storage.Delete(id)
		return err
	}
	return nil
}

func(s *OmaFileService) uploadToObjectStorage(ctx context.Context, input UploadFileInput, oma models.Omafile) error {
	file, err := input.File.Open()
	defer file.Close()
	if err != nil {
		return fmt.Errorf("open oma file: %w", err)
	}
	err = s.r2.UploadFileToR2(ctx, oma.OMAKey, file)
	if err != nil {
		return fmt.Errorf("open oma file: %w", err)
	}
	if input.Image == nil {
		return nil
	}

	img, err := input.Image.Open()
	defer img.Close()
	if err != nil {
		return ErrImgUpload
	}
	err = s.r2.UploadFileToR2(ctx, oma.ImgKey, img)
	if err != nil {
		return ErrImgUpload
	}
	return nil
}

func (s *OmaFileService) SearchOmaFile(ctx context.Context, brand string, model string) []models.Omafile {
	brand = strings.ToLower(brand)
	model = strings.ToLower(model)

	var files []models.Omafile = nil

	if brand != "" && model != "" {
		files = s.storage.GetByBrandAndModel(brand, model)
	} else if brand != "" {
		files = s.storage.GetByBrand(brand)
	} else if model != "" {
		files = s.storage.GetByModel(model)
	}

	for i := range files {
		if files[i].ImgKey != "" {
			url, err := s.r2.GeneratePresignedURLForImg(ctx, files[i].ImgKey)
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			files[i].ImgURL = url
		}
	}
	return files
}

func (s *OmaFileService) DownloadFile(ctx context.Context, id string) (DownloadFileOutput, error) {
	var output DownloadFileOutput
	
	oma := s.storage.GetById(id)
	object, err := s.r2.DownloadFileFromR2(ctx, oma.OMAKey)
	if err != nil {
		slog.Error(err.Error())
		return output, err
	}

	output.Oma = oma
	output.File = object

	return output, nil
}