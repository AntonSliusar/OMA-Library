package main

import (
	"context"
	"fmt"
	"oma-library/internal/config"
	"oma-library/internal/handlers"
	"oma-library/pkg/logger"
	"oma-library/pkg/storage"

	_ "github.com/lib/pq"
)

func main() {
	logger.Init()
	ctx := context.Background()
	fmt.Println("Context created")
	cfg := config.SetConfig()
	dbStorage, err := storage.NewStorage(cfg)
	if err != nil {
		panic(err)
	}
	r2, err := storage.NewR2Client(ctx, cfg.R2)
	if err != nil {
		panic(err)
	}
	fmt.Println(dbStorage, r2)
	handlers.RunWeb(dbStorage, r2)
}
