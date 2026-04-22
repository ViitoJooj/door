package jwtTokens

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func refreshSecret() string {
	secret := os.Getenv("JWT_REFRESH_TOKEN_SECRET")
	if secret == "" {
		panic("JWT_REFRESH_TOKEN_SECRET não definida")
	}
	return secret
}

func accessSecret() string {
	secret := os.Getenv("JWT_ACCESS_TOKEN_SECRET")
	if secret == "" {
		panic("JWT_ACCESS_TOKEN_SECRET não definida")
	}
	return secret
}

func GenerateRefreshToken(userID int) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"iat":     now.Unix(),
		"nbf":     now.Unix(),
		"exp":     now.Add(30 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(refreshSecret()))
}

func ValidateRefreshToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(refreshSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	return token, nil
}

func GenerateAccessToken(userID int) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "access",
		"iat":     now.Unix(),
		"nbf":     now.Unix(),
		"exp":     now.Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(accessSecret()))
}

func ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(accessSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return nil, errors.New("invalid token type")
	}

	return token, nil
}
