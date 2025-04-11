package migrator

import (
	"fmt"

	"github.com/AntonyCarl/OMA-Library/pkg/logger"
	"github.com/AntonyCarl/OMA-Library/pkg/storage"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
)

func Migrate(storage *storage.Storage) {
	driver, err := postgres.WithInstance(storage.GetDB(), &postgres.Config{})
	if err != nil {
		logger.Logger.Fatal("Не вдалося ініціалізувати драйвер PostgreSQL:", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://../db/migrations",
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
