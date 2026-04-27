package services

import (
	"database/sql"
	"errors"
	"log"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/repository"
)

const (
	defaultRequestsPerSecond = 1.0
	defaultBurst             = 5
	defaultProgressive       = false
)

type RateLimitService struct {
	RateLimitRepo repository.RateLimitRepository
}

func NewRateLimitService(rateLimitRepo repository.RateLimitRepository) *RateLimitService {
	return &RateLimitService{
		RateLimitRepo: rateLimitRepo,
	}
}

func (s RateLimitService) Get() (*domain.RateLimitSettings, error) {
	settings, err := s.RateLimitRepo.GetRateLimitSettings()
	if err == nil {
		return settings, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		defaultSettings, createErr := domain.NewRateLimitSettings(defaultRequestsPerSecond, defaultBurst, defaultProgressive)
		if createErr != nil {
			return nil, createErr
		}

		if upsertErr := s.RateLimitRepo.UpsertRateLimitSettings(defaultSettings); upsertErr != nil {
			log.Println(upsertErr)
			return nil, errors.New("internal error")
		}

		settings, readErr := s.RateLimitRepo.GetRateLimitSettings()
		if readErr != nil {
			log.Println(readErr)
			return nil, errors.New("internal error")
		}
		return settings, nil
	}

	log.Println(err)
	return nil, errors.New("internal error")
}

func (s RateLimitService) Update(requestsPerSecond float64, burst int, progressive bool) (*domain.RateLimitSettings, error) {
	settings, err := domain.NewRateLimitSettings(requestsPerSecond, burst, progressive)
	if err != nil {
		return nil, err
	}

	if err := s.RateLimitRepo.UpsertRateLimitSettings(settings); err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}

	updated, err := s.RateLimitRepo.GetRateLimitSettings()
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}

	return updated, nil
}
