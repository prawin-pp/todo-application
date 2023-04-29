package mock

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/parwin-pp/todo-application/internal/model"
)

type AuthDatabase struct {
	GetUserByUsernameFn func(ctx context.Context, username string) (*model.User, error)
	GetUserFn           func(ctx context.Context, userID string) (*model.User, error)
}

func (db *AuthDatabase) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	return db.GetUserByUsernameFn(ctx, username)
}

func (db *AuthDatabase) GetUser(ctx context.Context, userID string) (*model.User, error) {
	return db.GetUserFn(ctx, userID)
}

type AuthEncryptor struct {
	HashFn            func(str string) (string, error)
	CompareHashFn     func(hashedStr string, compareStr string) error
	SignAuthTokenFn   func(subject string, claims map[string]interface{}) (string, error)
	VerifyAuthTokenFn func(tokenStr string) (*jwt.Token, *jwt.MapClaims, error)
}

func (en *AuthEncryptor) Hash(str string) (string, error) {
	return en.HashFn(str)
}

func (en *AuthEncryptor) CompareHash(hashedStr string, compareStr string) error {
	return en.CompareHashFn(hashedStr, compareStr)
}

func (en *AuthEncryptor) SignAuthToken(subject string, claims map[string]interface{}) (string, error) {
	return en.SignAuthTokenFn(subject, claims)
}

func (en *AuthEncryptor) VerifyAuthToken(tokenStr string) (*jwt.Token, *jwt.MapClaims, error) {
	return en.VerifyAuthTokenFn(tokenStr)
}
