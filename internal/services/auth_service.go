package services

import (
	"database/sql"
	"errors"
	"strconv"

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

	newUser, err := domain.NewUser(
		user.Username,
		user.Email,
		user.Password,
	)
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

func (s *AuthService) Login(username string, email string, password string, ip string) (*domain.User, string, error) {
	if len(username) > 1 {
		user, err := s.userRepo.FindUserByUsername(username)
		if err != nil {
			s.logger.Error("failed to find user by username / error: " + err.Error())
			return nil, "", errors.New("internal error")
		}
		if user == nil {
			s.logger.Warn("login attempt with unknown username / username: " + username)
			return nil, "", errors.New("invalid credentials")
		}

		if !cryptography.CheckPasswordHash(password, user.Password) {
			s.logger.Warn("login attempt with wrong password / username: " + username)
			return nil, "", errors.New("invalid credentials")
		}

		token, err := jwtTokens.GenerateToken(user.ID)
		if err != nil {
			s.logger.Error("failed to generate token / error: " + err.Error())
			return nil, "", errors.New("internal error")
		}

		s.logger.Info("user logged in / user_id: " + strconv.Itoa(user.ID) + " | ip: " + ip)
		return user, token, nil
	} else {
		user, err := s.userRepo.FindUserByEmail(email)
		if err != nil {
			s.logger.Error("failed to find user by email / error: " + err.Error())
			return nil, "", errors.New("internal error")
		}
		if user == nil {
			s.logger.Warn("login attempt with unknown email  email: " + email)
			return nil, "", errors.New("invalid credentials")
		}

		if !cryptography.CheckPasswordHash(password, user.Password) {
			s.logger.Warn("login attempt with wrong password / email: " + email)
			return nil, "", errors.New("invalid credentials")
		}

		token, err := jwtTokens.GenerateToken(user.ID)
		if err != nil {
			s.logger.Error("failed to generate token / error: " + err.Error())
			return nil, "", errors.New("internal error")
		}

		s.logger.Info("user logged in / user_id: " + strconv.Itoa(user.ID) + " | ip: " + ip)
		return user, token, nil
	}
}

func (s *AuthService) Token(tokenString string) (*domain.User, error) {
	token, err := jwtTokens.ValidateToken(tokenString)
	if err != nil {
		s.logger.Warn("invalid or expired token / error: " + err.Error())
		return nil, errors.New("internal error")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Warn("failed to parse token claims")
		return nil, errors.New("invalid token")
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		s.logger.Warn("token missing user_id claim")
		return nil, errors.New("invalid token")
	}

	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		s.logger.Error("failed to find user by id / error: " + err.Error())
		return nil, errors.New("internal error")
	}
	if user == nil {
		s.logger.Warn("token references non-existing user")
		return nil, errors.New("invalid token")
	}
	return user, nil
}
