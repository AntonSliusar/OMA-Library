package repository

import (
	"io"
	"mime/multipart"
	"oma-library/pkg/logger"
	"os"
)

const path = "C:/GoLearn/FileRep/" // add to config

func SaveFile(file multipart.File, fileName string) string {
	dst, err := os.Create(path + fileName)
	if err != nil {
		logger.Logger.Error(err)
	}
	defer dst.Close()

	io.Copy(dst, file)
	return path + fileName
}
