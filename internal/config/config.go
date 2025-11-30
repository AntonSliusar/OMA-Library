package config

import (
	"fmt"
	"oma-library/pkg/logger"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB PostgresConfig
	R2 R2Config
}

type R2Config struct {
	Endpoint  string `envconfig:"ENDPOINT"`
	AccessKey string `envconfig:"ACCESS_KEY"`
	SecretKey string `envconfig:"SECRET_KEY"`	
	Bucket    string `envconfig:"BUCKET"`
}

type PostgresConfig struct {
	Host     string `envconfig:"HOST"`
	Port     int	`envconfig:"PORT"`
	User     string	`envconfig:"USER"`
	Password string	`envconfig:"PASSWORD"`
	DBName   string	`envconfig:"DBNAME"`
	SSLMode  string	`envconfig:"SSLMODE"`
}

func SetConfig() *Config {
	var cfg Config

	err := envconfig.Process("db", &cfg.DB)
	if err != nil {
		fmt.Println(err)
		logger.Logger.Fatal(err)
	}

	err = envconfig.Process("r2", &cfg.R2)
	if err != nil {
		fmt.Println(err)
		logger.Logger.Fatal(err)
	}

	return &cfg
}

// config for admin DB
