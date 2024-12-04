package config

import "github.com/joho/godotenv"

// PGConfigI предоставляет контракт для получения адреса базы данных
type PGConfigI interface {
	DSN() string
}

// GRPCConfigI предоставляет контракт для получения адреса GRPC-сервера
type GRPCConfigI interface {
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
