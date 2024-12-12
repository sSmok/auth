package config

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	refreshSecretKey  = "REFRESH_SECRET_KEY"
	accessSecretKey   = "ACCESS_SECRET_KEY"
	refreshExpiration = "REFRESH_EXPIRATION"
	accessExpiration  = "ACCESS_EXPIRATION"
)

// TokenConfigI - интерфейс конфигурации токена
type TokenConfigI interface {
	RefreshTokenSecretKey() string
	AccessTokenSecretKey() string
	RefreshTokenExpiration() time.Duration
	AccessTokenExpiration() time.Duration
}

type token struct {
	refreshSecretKey  string
	accessSecretKey   string
	refreshExpiration time.Duration
	accessExpiration  time.Duration
}

// NewTokenConfig - конструктор конфигурации токена
func NewTokenConfig() (TokenConfigI, error) {
	refreshKey := os.Getenv(refreshSecretKey)
	if len(refreshKey) == 0 {
		return nil, errors.New("refreshSecretKey not found")
	}
	accessKey := os.Getenv(accessSecretKey)
	if len(accessKey) == 0 {
		return nil, errors.New("accessSecretKey not found")
	}
	refreshExp := os.Getenv(refreshExpiration)
	if len(refreshExp) == 0 {
		return nil, errors.New("refreshExpiration not found")
	}
	refreshExpInt, err := strconv.Atoi(refreshExp)
	if err != nil {
		return nil, err
	}
	accessExp := os.Getenv(accessExpiration)
	if len(accessExp) == 0 {
		return nil, errors.New("accessExpiration not found")
	}
	accessExpInt, err := strconv.Atoi(accessExp)
	if err != nil {
		return nil, err
	}

	t := &token{
		refreshSecretKey:  refreshKey,
		accessSecretKey:   accessKey,
		refreshExpiration: time.Duration(refreshExpInt) * time.Minute,
		accessExpiration:  time.Duration(accessExpInt) * time.Minute,
	}

	return t, nil
}

func (t *token) RefreshTokenSecretKey() string {
	return t.refreshSecretKey
}

func (t *token) AccessTokenSecretKey() string {
	return t.accessSecretKey
}

func (t *token) RefreshTokenExpiration() time.Duration {
	return t.refreshExpiration
}

func (t *token) AccessTokenExpiration() time.Duration {
	return t.accessExpiration
}
