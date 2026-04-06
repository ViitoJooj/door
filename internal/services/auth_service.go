package services

import (
	"errors"
	"log"

	"github.com/ViitoJooj/door/internal/domain"
	"github.com/ViitoJooj/door/internal/repository"
	"github.com/ViitoJooj/door/pkg/cryptography"
	"github.com/ViitoJooj/door/pkg/jwtTokens"
	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(user *domain.User) (*domain.User, error) {
	existing, err := s.userRepo.FindUserByEmail(user.Email)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	if existing != nil {
		log.Println("User already exists")
		return nil, errors.New("invalid credentials")
	}

	newUser, err := domain.NewUser(
		user.Username,
		user.Email,
		user.Password,
	)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}

	if err := s.userRepo.CreateUser(newUser); err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}

	return newUser, nil
}

func (s *AuthService) Login(username string, email string, password string) (*domain.User, string, error) {
	if len(username) > 1 {
		user, err := s.userRepo.FindUserByUsername(username)
		if err != nil {
			log.Println(err)
			return nil, "", errors.New("internal error")
		}
		if user == nil {
			log.Println("User not found")
			return nil, "", errors.New("invalid credentials")
		}

		if !cryptography.CheckPasswordHash(password, user.Password) {
			log.Println("Invalid password.")
			return nil, "", errors.New("invalid credentials")
		}

		token, err := jwtTokens.GenerateToken(user.ID)
		if err != nil {
			log.Println(err)
			return nil, "", errors.New("internal error")
		}

		return user, token, nil
	} else {
		user, err := s.userRepo.FindUserByEmail(email)
		if err != nil {
			log.Println(err)
			return nil, "", errors.New("internal error")
		}
		if user == nil {
			log.Println("User not found")
			return nil, "", errors.New("invalid credentials")
		}

		if !cryptography.CheckPasswordHash(password, user.Password) {
			log.Println("Invalid password.")
			return nil, "", errors.New("invalid credentials")
		}

		token, err := jwtTokens.GenerateToken(user.ID)
		if err != nil {
			log.Println(err)
			return nil, "", errors.New("internal error")
		}

		return user, token, nil
	}
}

func (s *AuthService) Token(tokenString string) (*domain.User, error) {
	token, err := jwtTokens.ValidateToken(tokenString)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("invalid token claims")
		return nil, errors.New("invalid token")
	}
	userID, ok := claims["user_id"].(int)
	if !ok {
		log.Println("invalid token claims")
		return nil, errors.New("invalid token")
	}

	user, err := s.userRepo.FindUserByID(userID)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	if user == nil {
		log.Println("User not exists.")
		return nil, errors.New("invalid token")
	}
	return user, nil
}
