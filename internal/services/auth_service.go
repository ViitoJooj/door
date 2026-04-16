package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ViitoJooj/door/internal/domain"
	"github.com/ViitoJooj/door/internal/repository"
	"github.com/ViitoJooj/door/pkg/cryptography"
	"github.com/ViitoJooj/door/pkg/jwtTokens"
	"github.com/ViitoJooj/door/pkg/logger"
	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	userRepo repository.UserRepository
	logger   *logger.Logger
}

func NewAuthService(userRepo repository.UserRepository, log *logger.Logger) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		logger:   log,
	}
}

func (s *AuthService) Register(user *domain.User) (*domain.User, error) {
	existing, err := s.userRepo.FindUserByEmail(user.Email)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Error("failed to find user by email / error: " + err.Error())
		return nil, errors.New("internal error")
	}
	if existing != nil {
		s.logger.Warn("register attempt with existing email / email: " + user.Email)
		return nil, errors.New("invalid credentials")
	}

	newUser, err := domain.NewUser(user.Username, user.Email, user.Password)
	if err != nil {
		s.logger.Error("failed to create user domain / error: " + err.Error())
		return nil, errors.New("internal error")
	}

	hashedPassword, err := cryptography.HashPassword(newUser.Password)
	if err != nil {
		s.logger.Error("failed to hash password / error: " + err.Error())
		return nil, errors.New("internal error")
	}

	newUser.Password = hashedPassword

	if err := s.userRepo.CreateUser(newUser); err != nil {
		s.logger.Error("failed to insert user / error: " + err.Error())
		return nil, errors.New("internal error")
	}

	return newUser, nil
}

func (s *AuthService) Login(username string, email string, password string, ip string) (*domain.User, string, string, error) {
	if len(username) > 1 {
		user, err := s.userRepo.FindUserByUsername(username)
		if err != nil {
			s.logger.Error("failed to find user by username / error: " + err.Error())
			return nil, "", "", errors.New("internal error")
		}
		if user == nil {
			s.logger.Warn("login attempt with unknown username / username: " + username)
			return nil, "", "", errors.New("invalid credentials")
		}

		if !cryptography.CheckPasswordHash(password, user.Password) {
			s.logger.Warn("login attempt with wrong password / username: " + username)
			return nil, "", "", errors.New("invalid credentials")
		}

		accessToken, err := jwtTokens.GenerateAccessToken(user.ID)
		if err != nil {
			return nil, "", "", errors.New("internal error")
		}

		refreshToken, err := jwtTokens.GenerateRefreshToken(user.ID)
		if err != nil {
			return nil, "", "", errors.New("internal error")
		}

		textLog := fmt.Sprintf("user logged in / user_id: %d | ip: %s", user.ID, ip)
		s.logger.Info(textLog)
		return user, accessToken, refreshToken, nil
	} else {
		user, err := s.userRepo.FindUserByEmail(email)
		if err != nil {
			s.logger.Error("failed to find user by email / error: " + err.Error())
			return nil, "", "", errors.New("internal error")
		}
		if user == nil {
			s.logger.Warn("login attempt with unknown email  email: " + email)
			return nil, "", "", errors.New("invalid credentials")
		}

		if !cryptography.CheckPasswordHash(password, user.Password) {
			s.logger.Warn("login attempt with wrong password / email: " + email)
			return nil, "", "", errors.New("invalid credentials")
		}

		accessToken, err := jwtTokens.GenerateAccessToken(user.ID)
		if err != nil {
			return nil, "", "", errors.New("internal error")
		}

		refreshToken, err := jwtTokens.GenerateRefreshToken(user.ID)
		if err != nil {
			return nil, "", "", errors.New("internal error")
		}

		textLog := fmt.Sprintf("user logged in / user_id: %d | ip: %s", user.ID, ip)
		s.logger.Info(textLog)
		return user, accessToken, refreshToken, nil
	}
}

func (s *AuthService) Token(tokenString string, isRefresh bool) (*domain.User, error) {
	var token *jwt.Token
	var err error

	if isRefresh {
		token, err = jwtTokens.ValidateRefreshToken(tokenString)
	} else {
		token, err = jwtTokens.ValidateAccessToken(tokenString)
	}

	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token")
	}

	userID := int(userIDFloat)

	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("internal error")
	}
	if user == nil {
		return nil, errors.New("invalid token")
	}

	return user, nil
}
