package storage

import (
	"database/sql"
	"fmt"
	"log/slog"
	"oma-library/internal/config"
	"oma-library/internal/models"
)

type Storage struct {
	db *sql.DB

}

func NewStorage(cfg *config.Config) (*Storage, error) {
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.DBName,
			cfg.DB.SSLMode))
	if err != nil {
		slog.Error("database NewStorage:", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		slog.Error("database NewStorage:", err.Error())
	}
	return &Storage{db: db}, nil
}

func (storage *Storage) Create(o models.Omafile) (int64, error) {
	var id int64
	err := storage.db.QueryRow("INSERT INTO files (brand, model, omakey, imgkey) VALUES ($1, $2, $3, $4) RETURNING id",
		o.Brand, o.Model, o.OMAKey, o.ImgKey).Scan(&id)

	if err != nil {
		slog.Error("database Create:", err.Error())
	}
	return id, err
}

func (storage *Storage) GetAll() ([]models.Omafile, error) {
	rows, err := storage.db.Query("SELECT * FROM files")
	if err != nil {
		slog.Error("database GetAll:", err.Error())
	}
	defer rows.Close()
	forms := make([]models.Omafile, 0)
	for rows.Next() {
		form := models.Omafile{}
		err := rows.Scan(&form.Id, &form.Brand, &form.Model, &form.OMAKey, &form.ImgKey)
		if err != nil {
			slog.Error("database GetAll:", err.Error())
		}
		forms = append(forms, form)
	}
	return forms, nil
}

func (storage *Storage) GetById(id string) models.Omafile {
	rows, err := storage.db.Query("SELECT * FROM files WHERE id = $1", id)
	if err != nil {
		slog.Error("database GetById:", err.Error())
	}
	defer rows.Close()

	form := models.Omafile{}
	for rows.Next() {
		err := rows.Scan(&form.Id, &form.Brand, &form.Model, &form.OMAKey, &form.ImgKey)
		if err != nil {
			slog.Error("database GetById:", err.Error())
		}
	}
	return form
}

func (storage *Storage) GetByBrand(brand string) []models.Omafile {
	rows, err := storage.db.Query("SELECT * FROM files WHERE brand = $1", brand)
	if err != nil {
		slog.Error("database GetByBrand:",err.Error())
	}
	defer rows.Close()

	forms := make([]models.Omafile, 0)
	for rows.Next() {
		form := models.Omafile{}
		err := rows.Scan(&form.Id, &form.Brand, &form.Model, &form.OMAKey, &form.ImgKey)
		if err != nil {
			slog.Error("database GetByBrand:",err.Error())
		}
		forms = append(forms, form)
	}
	return forms
}

func (storage *Storage) GetByBrandAndModel(brand string, model string) []models.Omafile {
	rows, err := storage.db.Query("SELECT * FROM files WHERE brand = $1 AND model = $2", brand, model)
	if err != nil {
		slog.Error("database GetBybrandAndModel:",err.Error())
	}
	defer rows.Close()

	forms := make([]models.Omafile, 0)
	for rows.Next() {
		form := models.Omafile{}
		err := rows.Scan(&form.Id, &form.Brand, &form.Model, &form.OMAKey, &form.ImgKey)
		if err != nil {
			slog.Error("database GetByBrandAndModel:",err.Error())
		}
		forms = append(forms, form)
	}
	return forms
}

func (storage *Storage) GetByModel(model string) []models.Omafile {
	rows, err := storage.db.Query("SELECT * FROM files WHERE model = $1", model)
	if err != nil {
		slog.Error("database GetByModel:", err.Error())
	}
	defer rows.Close()

	forms := make([]models.Omafile, 0)
	for rows.Next() {
		form := models.Omafile{}
		err := rows.Scan(&form.Id, &form.Brand, &form.Model, &form.OMAKey, &form.ImgKey)
		if err != nil {
			slog.Error("database GetByModel:", err.Error())
		}
		forms = append(forms, form)
	}
	return forms
}

func (storage *Storage) Delete(id int64) error {
	_, err := storage.db.Exec("DELETE FROM files WHERE id = $1", id)
	if err != nil {
		slog.Error("database Delete func:", err.Error())
		return err
	}
	return err
}

func (storage *Storage) GetDB() *sql.DB {
	return storage.db
}
