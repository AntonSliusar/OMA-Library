package migrator

import (
	"fmt"

	"oma-library/pkg/logger"

	"oma-library/pkg/storage"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(storage *storage.Storage) {
	driver, err := postgres.WithInstance(storage.GetDB(), &postgres.Config{})
	if err != nil {
		logger.Logger.Fatal("Не вдалося ініціалізувати драйвер PostgreSQL:", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://pkg/db/migrations",
		"postgres", driver)
	if err != nil {
		logger.Logger.Fatal("Не вдалося створити мігратор:", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Logger.Fatal("Помилка застосування міграцій:", err)
	} else {
		fmt.Println("Міграції застосовано успішно або змін не знайдено.")
	}
}
