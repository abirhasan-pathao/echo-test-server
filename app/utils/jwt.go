package utils

import (
	"echo-server/config"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: " + err.Error())
	}
}

func GenerateJWT() (string, error) {
	if cfg.Environment.JWTSecret == "" {
		return "", fmt.Errorf("jwt secret is not set in configuration")
	}
	jwtSecret := []byte(cfg.Environment.JWTSecret)
	claims := jwt.MapClaims{
		"exp": time.Now().Add(1 * time.Minute).Unix(), // 1m expiry
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (bool, error) {
	jwtSecret := []byte(cfg.Environment.JWTSecret)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		return false, err
	}
	return token.Valid, nil
}