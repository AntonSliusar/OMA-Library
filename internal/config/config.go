package config

import (
	"log/slog"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB PostgresConfig
	R2 R2Config
	MCP MCPConfig
	AUTH Auth
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

type MCPConfig struct {
	URL string `envconfig:"URL"`
	Port string `envconfig:"PORT"`
}

type Auth struct {
	JWT string `envconfig:"TOKEN"`
}

func SetConfig() *Config {
	var cfg Config

	err := envconfig.Process("db", &cfg.DB)
	if err != nil {
		slog.Error(err.Error())
	}

	err = envconfig.Process("r2", &cfg.R2)
	if err != nil {
		slog.Error(err.Error())
	}

	err = envconfig.Process("mcp", &cfg.MCP)
	if err != nil {
		slog.Error(err.Error())
	}

	err = envconfig.Process("jwt", &cfg.AUTH)
	if err != nil {
		slog.Error(err.Error())
	}

	return &cfg
}

// config for admin DB
