package main

import (
	"context"
	"oma-library/internal/config"
	"oma-library/internal/handlers"
	"oma-library/pkg/db/migrator"
	"oma-library/pkg/logger"
	"oma-library/pkg/storage"

	_ "github.com/lib/pq"
)

func main() {
	//TODO: add admin table to migration files
	ctx := context.Background()
	cfg := config.SetConfig()
	logger.Init()
	dbStorage, _ := storage.NewStorage(cfg)
	r2, _ := storage.NewR2Client(ctx, cfg.R2)
	migrator.Migrate(dbStorage)
	handlers.RunWeb(dbStorage, r2)
}
