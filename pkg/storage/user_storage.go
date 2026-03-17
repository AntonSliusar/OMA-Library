package storage

import (
	"log/slog"
	"oma-library/internal/models"
)

func (storage *Storage) AddUser(req models.SignUpRequest) error {
	_, err := storage.db.Exec("INSERT INTO admins (username, email, password) VALUES ($1, $2, $3)",
		req.Username, req.Email, req.Password)
	if err != nil {
		slog.Error(err.Error())
	}
	return err
}

func (storage *Storage) CheckExist(email string) (bool, error) {
	var exists bool
	err := storage.db.QueryRow("SELECT EXISTS(SELECT * FROM admins WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		slog.Error(err.Error())
		return false, err
	}
	return exists, nil
}

func (storage *Storage) GetByEmail(email string) (models.User, error) {
	row := storage.db.QueryRow("SELECT id, username, email, password FROM admins WHERE email=$1", email)
	admin := models.User{}
	err := row.Scan(&admin.ID, &admin.Username, &admin.Email, &admin.Password)
	if err != nil {
		slog.Error(err.Error())
	}
	return admin, err
}
