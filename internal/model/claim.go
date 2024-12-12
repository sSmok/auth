package model

import "github.com/dgrijalva/jwt-go"

// UserClaims - данные токена
type UserClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	Role  int32  `json:"role"`
}
