package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthEncryption struct {
	SigningMethod jwt.SigningMethod
	SigningKey    interface{}
	TTL           time.Duration
}

func NewAuthEncryption(method string, secret interface{}, ttl string) *AuthEncryption {
	duration, err := time.ParseDuration(ttl)
	if err != nil {
		log.Fatalf("failed to parse ttl: %v", err)
	}
	return &AuthEncryption{
		SigningMethod: jwt.GetSigningMethod(method),
		SigningKey:    secret,
		TTL:           duration,
	}
}

func (ae *AuthEncryption) Hash(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (ae *AuthEncryption) CompareHash(hashedStr string, compareStr string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(compareStr))
}

func (ae *AuthEncryption) SignAuthToken(subject string, claims map[string]interface{}) (string, error) {
	authClaims := ae.createAuthClaims(subject, claims)
	token := jwt.NewWithClaims(ae.SigningMethod, authClaims)
	return token.SignedString(ae.SigningKey)
}

func (ae *AuthEncryption) createAuthClaims(subject string, data map[string]interface{}) *jwt.MapClaims {
	claims := jwt.MapClaims{}
	for key, value := range data {
		claims[key] = value
	}
	claims["sub"] = subject
	claims["exp"] = time.Now().Add(ae.TTL).Unix()
	return &claims
}

func (ae *AuthEncryption) VerifyAuthToken(tokenStr string) (*jwt.Token, *jwt.MapClaims, error) {
	authClaims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, authClaims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != ae.SigningMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return ae.SigningKey, nil
	})
	if err != nil {
		return nil, nil, err
	}
	return token, authClaims, nil
}
