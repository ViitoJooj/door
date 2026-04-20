package services

import (
	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/repository"
)

type DotEnvService struct {
	dotEnvRepo repository.DotEnvRepository
}

func NewDotEnvService(DotEnvRepo repository.DotEnvRepository) *DotEnvService {
	return &DotEnvService{
		dotEnvRepo: DotEnvRepo,
	}
}

func (s *DotEnvService) GetAll() ([]*domain.Env, error) {
	return s.dotEnvRepo.GetAllVars()
}

func (s *DotEnvService) GetVar(id int) (*domain.Env, error) {
	return s.dotEnvRepo.FindVar(id)
}

func (s *DotEnvService) ChangeVar(env domain.Env) error {
	return s.dotEnvRepo.ChangeVar(&env)
}
