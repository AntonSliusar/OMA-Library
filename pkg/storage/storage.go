package storage

import (
	"database/sql"
	"fmt"
	"oma-library/internal/config"
	"oma-library/internal/models"

	"oma-library/pkg/logger"
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
		logger.Logger.Fatal(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	return &Storage{db: db}, nil
}

func (storage *Storage) Create(o models.Omafile) error {
	_, err := storage.db.Exec("INSERT INTO files (brand, model, omakey, imgkey) VALUES ($1, $2, $3, $4)",
		o.Brand, o.Model, o.OMAKey, o.ImgKey)

	if err != nil {
		logger.Logger.Error(err)
	}
	return err
}

func (storage *Storage) GetById(id string) models.Omafile {
	rows, err := storage.db.Query("SELECT * FROM files WHERE id = $1", id)
	if err != nil {
		logger.Logger.Error(err)
	}
	defer rows.Close()

	form := models.Omafile{}
	for rows.Next() {
		err := rows.Scan(&form.Id, &form.Brand, &form.Model, &form.OMAKey, &form.ImgKey)
		if err != nil {
			logger.Logger.Error(err)
		}
	}
	return form
}

func (storage *Storage) GetByBrand(brand string) []models.Omafile {
	rows, err := storage.db.Query("SELECT * FROM files WHERE brand = $1", brand)
	if err != nil {
		logger.Logger.Error(err)
	}
	defer rows.Close()

	forms := make([]models.Omafile, 0)
	for rows.Next() {
		form := models.Omafile{}
		err := rows.Scan(&form.Id, &form.Brand, &form.Model, &form.OMAKey, &form.ImgKey)
		if err != nil {
			logger.Logger.Error(err)
		}
		forms = append(forms, form)
	}
	return forms
}

func (storage *Storage) GetByBrandAndModel(brand string, model string) []models.Omafile {
	rows, err := storage.db.Query("SELECT * FROM files WHERE brand = $1 AND model = $2", brand, model)
	if err != nil {
		logger.Logger.Error(err)
	}
	defer rows.Close()

	forms := make([]models.Omafile, 0)
	for rows.Next() {
		form := models.Omafile{}
		err := rows.Scan(&form.Id, &form.Brand, &form.Model, &form.OMAKey, &form.ImgKey)
		if err != nil {
			logger.Logger.Error(err)
		}
		forms = append(forms, form)
	}
	return forms
}

func (storage *Storage) GetByModel(model string) []models.Omafile {
	rows, err := storage.db.Query("SELECT * FROM files WHERE model = $1", model)
	if err != nil {
		logger.Logger.Error(err)
	}
	defer rows.Close()

	forms := make([]models.Omafile, 0)
	for rows.Next() {
		form := models.Omafile{}
		err := rows.Scan(&form.Id, &form.Brand, &form.Model, &form.OMAKey, &form.ImgKey)
		if err != nil {
			logger.Logger.Error(err)
		}
		forms = append(forms, form)
	}
	return forms
}

func (storage *Storage) GetDB() *sql.DB {
	return storage.db
}
