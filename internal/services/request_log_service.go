package services

import (
	"log"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/repository"
)

type RequestLogService struct {
	LogRepo repository.RequestLogRepository
}

func NewRequestLogService(logRepo repository.RequestLogRepository) *RequestLogService {
	return &RequestLogService{
		LogRepo: logRepo,
	}
}

func (s RequestLogService) GetAll() ([]*domain.RequestLog, error) {
	logs, err := s.LogRepo.ListRequestLogs()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return logs, nil
}
