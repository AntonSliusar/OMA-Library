package storage

import (
	"oma-library/internal/models"

	"oma-library/pkg/logger"
)

func (storage *Storage) AddAdmin(req models.SignUpRequest) error {
	_, err := storage.db.Exec("INSERT INTO admins (username, email, password) VALUES ($1, $2, $3)",
		req.Username, req.Email, req.Password)
	if err != nil {
		logger.Logger.Error(err)
	}
	return err
}

func (storage *Storage) CheckExist(email string) bool {
	var exists bool
	err := storage.db.QueryRow("SELECT EXISTS(SELECT * FROM admins WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		logger.Logger.Error(err)
	}
	return exists
}

func (storage *Storage) GetByEmail(email string) (models.Admin, error) {
	row := storage.db.QueryRow("SELECT id, username, email, password FROM admins WHERE email=$1", email)
	admin := models.Admin{}
	err := row.Scan(&admin.ID, &admin.Username, &admin.Email, &admin.Password)
	if err != nil {
		logger.Logger.Error(err)
	}
	return admin, err
}
