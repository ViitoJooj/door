package services

import (
	"database/sql"
	"errors"
	"log"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/repository"
)

const (
	defaultAllowedProtocol = domain.ProtocolModeBoth
	defaultApplyScope      = domain.ConfigScopeAll
)

type ProtocolSettingsService struct {
	Repo repository.ProtocolSettingsRepository
}

func NewProtocolSettingsService(repo repository.ProtocolSettingsRepository) *ProtocolSettingsService {
	return &ProtocolSettingsService{
		Repo: repo,
	}
}

func (s ProtocolSettingsService) Get() (*domain.ProtocolSettings, error) {
	settings, err := s.Repo.GetProtocolSettings()
	if err == nil {
		return settings, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		defaultSettings, createErr := domain.NewProtocolSettings(defaultAllowedProtocol, defaultApplyScope)
		if createErr != nil {
			return nil, createErr
		}

		if upsertErr := s.Repo.UpsertProtocolSettings(defaultSettings); upsertErr != nil {
			log.Println(upsertErr)
			return nil, errors.New("internal error")
		}

		settings, readErr := s.Repo.GetProtocolSettings()
		if readErr != nil {
			log.Println(readErr)
			return nil, errors.New("internal error")
		}
		return settings, nil
	}

	log.Println(err)
	return nil, errors.New("internal error")
}

func (s ProtocolSettingsService) Update(allowedProtocol string, applyScope string) (*domain.ProtocolSettings, error) {
	settings, err := domain.NewProtocolSettings(allowedProtocol, applyScope)
	if err != nil {
		return nil, err
	}

	if err := s.Repo.UpsertProtocolSettings(settings); err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}

	updated, err := s.Repo.GetProtocolSettings()
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}

	return updated, nil
}
