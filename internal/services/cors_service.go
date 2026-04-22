package services

import (
	"errors"
	"log"
	"strings"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/repository"
)

type CorsService struct {
	CorsRepo repository.CorsRepository
	AuthRepo repository.UserRepository
}

func NewCorsService(corsRepo repository.CorsRepository, authRepo repository.UserRepository) *CorsService {
	return &CorsService{
		CorsRepo: corsRepo,
		AuthRepo: authRepo,
	}
}

func (s CorsService) Create(origin string, userID int) (*domain.Cors, *domain.User, error) {
	origin = strings.TrimSpace(origin)

	if origin == "" {
		return nil, nil, errors.New("origin is required")
	}

	user, err := s.AuthRepo.FindUserByID(userID)
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal error")
	} else if user == nil {
		return nil, nil, errors.New("internal error")
	}

	newCors := &domain.Cors{
		Origin: origin,
	}

	if err := s.CorsRepo.CreateCors(newCors); err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal error")
	}

	return newCors, user, nil
}

func (s CorsService) GetAll() ([]*domain.Cors, error) {
	return s.CorsRepo.FindAllCors()
}

func (s CorsService) GetByID(id int) (*domain.Cors, error) {
	cors, err := s.CorsRepo.FindCorsByID(id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	} else if cors == nil {
		return nil, errors.New("cors origin not found")
	}

	return cors, nil
}

func (s CorsService) Update(id int, origin string, userID int) (*domain.Cors, *domain.User, error) {
	origin = strings.TrimSpace(origin)

	if origin == "" {
		return nil, nil, errors.New("origin is required")
	}

	user, err := s.AuthRepo.FindUserByID(userID)
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal error")
	} else if user == nil {
		return nil, nil, errors.New("internal error")
	}

	exists, err := s.CorsRepo.FindCorsByID(id)
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal error")
	} else if exists == nil {
		return nil, nil, errors.New("cors origin not found")
	}

	exists.Origin = origin

	if err := s.CorsRepo.ChangeCors(exists); err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal error")
	}

	return exists, user, nil
}

func (s CorsService) DeleteByID(id int, userID int) (*domain.Cors, *domain.User, error) {
	user, err := s.AuthRepo.FindUserByID(userID)
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal error")
	} else if user == nil {
		return nil, nil, errors.New("internal error")
	}

	cors, err := s.CorsRepo.FindCorsByID(id)
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal error")
	} else if cors == nil {
		return nil, nil, errors.New("cors origin not found")
	}

	err = s.CorsRepo.DeleteCors(id)
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("internal error")
	}

	return cors, user, nil
}
