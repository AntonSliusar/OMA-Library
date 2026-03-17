package main

import (
	"context"
	"fmt"
	"log/slog"
	"oma-library/internal/config"
	"oma-library/internal/handlers"
	"oma-library/internal/server"
	"oma-library/internal/service"
	"oma-library/pkg/storage"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
	slog.SetDefault(logger)
	ctx := context.Background()
	cfg := config.SetConfig()
	dbStorage, err := storage.NewStorage(cfg)
	if err != nil {
		panic(err)
	}
	r2, err := storage.NewR2Client(ctx, cfg.R2)
	if err != nil {
		panic(err)
	}

	fileService := service.NewOmaFileService(dbStorage, r2)
	fileHandler := handlers.NewFileHandler(fileService)

	userService := service.NewUserService(dbStorage)
	userHandler := handlers.NewUserHandler(userService)

	fmt.Println(dbStorage, r2)
	// go handlers.RunMCP(dbStorage, r2, cfg)
	server.RunServer(cfg, fileHandler, userHandler)
}
 