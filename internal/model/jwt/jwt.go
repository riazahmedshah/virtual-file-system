package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID        uuid.UUID `json:"userId"`
	StorageLimit  int64     `json:"storageLimit"`
	FileSizeLimit int64     `json:"fileSizeLimit"`
	IsGuest       bool      `json:"isGuest"`
	jwt.RegisteredClaims
}
