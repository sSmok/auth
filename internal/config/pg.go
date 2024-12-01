package config

import (
	"errors"
	"os"
)

const dsnEnv = "DB_DSN"

type pgConfig struct {
	dsn string
}

// NewPGConfig достает из config-файла данные о базе данных: адрес БД
func NewPGConfig() (PGConfig, error) {
	dsn := os.Getenv(dsnEnv)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{dsn: dsn}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
