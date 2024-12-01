package config

import "github.com/joho/godotenv"

// PGConfig предоставляет контракт для получения адреса базы данных
type PGConfig interface {
	DSN() string
}

// GRPCConfig предоставляет контракт для получения адреса GRPC-сервера
type GRPCConfig interface {
	Address() string
}

// Load загружает указанный файл конфига, для последующей обработки
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
