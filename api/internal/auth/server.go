package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/parwin-pp/todo-application/internal/config"
	"github.com/parwin-pp/todo-application/internal/model"
)

type Database interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUser(ctx context.Context, userID string) (*model.User, error)
}

type Encrypter interface {
	Hash(str string) (string, error)
	CompareHash(hashedStr string, compareStr string) error
	SignAuthToken(subject string, claims map[string]interface{}) (string, error)
	VerifyAuthToken(tokenStr string) (*jwt.Token, *jwt.MapClaims, error)
}

type Server struct {
	db     Database
	en     Encrypter
	config config.AuthConfig
}

func NewServer(db Database, en Encrypter, config config.AuthConfig) *Server {
	return &Server{db: db, en: en, config: config}
}
