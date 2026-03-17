package repository

import (
	"io"
	"log/slog"
	"mime/multipart"
	"os"
)

const path = "C:/GoLearn/FileRep/" // add to config

func SaveFile(file multipart.File, fileName string) string {
	dst, err := os.Create(path + fileName)
	if err != nil {
		slog.Error(err.Error())
	}
	defer dst.Close()

	io.Copy(dst, file)
	return path + fileName
}
