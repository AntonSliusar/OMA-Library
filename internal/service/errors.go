package service

import "errors"

var (
	ErrInvalidFileFormat = errors.New("invalid file format, only .oma allowed")
	ErrFileNotFound      = errors.New("file not found")
	ErrImgUpload = errors.New("file uploaded without omage")
)