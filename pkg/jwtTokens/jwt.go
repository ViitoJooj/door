package jwtTokens

import (
	"errors"
	"time"

	"github.com/ViitoJooj/ward/pkg/dotenv"
	"github.com/golang-jwt/jwt/v4"
)

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
	return token.SignedString([]byte(dotenv.JwtRefreshTokenSecret))
}

func ValidateRefreshToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(dotenv.JwtRefreshTokenSecret), nil
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
	return token.SignedString([]byte(dotenv.JwtAccessTokenSecret))
}

func ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(dotenv.JwtAccessTokenSecret), nil
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
