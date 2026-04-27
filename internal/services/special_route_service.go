package services

import (
	"errors"
	"log"
	"strings"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/repository"
)

type SpecialRouteService struct {
	Repo     repository.SpecialRouteRepository
	AuthRepo repository.UserRepository
}

func NewSpecialRouteService(repo repository.SpecialRouteRepository, authRepo repository.UserRepository) *SpecialRouteService {
	return &SpecialRouteService{
		Repo:     repo,
		AuthRepo: authRepo,
	}
}

func (s SpecialRouteService) ensureUser(userID int) error {
	user, err := s.AuthRepo.FindUserByID(userID)
	if err != nil {
		log.Println(err)
		return errors.New("internal error")
	}
	if user == nil {
		return errors.New("internal error")
	}
	return nil
}

func isSpecialRouteUniqueErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}

func (s SpecialRouteService) GetByType(routeType string) ([]*domain.SpecialRouteRule, error) {
	normalizedType, err := domain.NormalizeSpecialRouteType(routeType)
	if err != nil {
		return nil, err
	}
	return s.Repo.ListSpecialRouteRules(normalizedType)
}

func (s SpecialRouteService) Create(routeType string, path string, maxDistinctRequests int, windowSeconds int, banSeconds int, enabled bool, userID int) (*domain.SpecialRouteRule, error) {
	if err := s.ensureUser(userID); err != nil {
		return nil, err
	}

	rule, err := domain.NewSpecialRouteRule(routeType, path, maxDistinctRequests, windowSeconds, banSeconds, enabled)
	if err != nil {
		return nil, err
	}
	rule.CreatedBy = userID
	rule.UpdatedBy = userID

	if err := s.Repo.CreateSpecialRouteRule(rule); err != nil {
		if isSpecialRouteUniqueErr(err) {
			return nil, errors.New("path already exists for this route_type")
		}
		log.Println(err)
		return nil, errors.New("internal error")
	}
	return rule, nil
}

func (s SpecialRouteService) Update(routeType string, id int, path string, maxDistinctRequests int, windowSeconds int, banSeconds int, enabled bool, userID int) (*domain.SpecialRouteRule, error) {
	if err := s.ensureUser(userID); err != nil {
		return nil, err
	}

	normalizedType, err := domain.NormalizeSpecialRouteType(routeType)
	if err != nil {
		return nil, err
	}

	rule, err := s.Repo.FindSpecialRouteRuleByID(id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	if rule == nil || rule.RouteType != normalizedType {
		return nil, errors.New("special route not found")
	}

	normalizedPath, err := domain.NormalizeSpecialRoutePath(path)
	if err != nil {
		return nil, err
	}
	if maxDistinctRequests <= 0 {
		return nil, errors.New("max_distinct_requests must be greater than 0")
	}
	if windowSeconds <= 0 {
		return nil, errors.New("window_seconds must be greater than 0")
	}
	if banSeconds <= 0 {
		return nil, errors.New("ban_seconds must be greater than 0")
	}

	rule.Path = normalizedPath
	rule.MaxDistinctRequests = maxDistinctRequests
	rule.WindowSeconds = windowSeconds
	rule.BanSeconds = banSeconds
	rule.Enabled = enabled
	rule.UpdatedBy = userID

	if err := s.Repo.UpdateSpecialRouteRule(rule); err != nil {
		if isSpecialRouteUniqueErr(err) {
			return nil, errors.New("path already exists for this route_type")
		}
		log.Println(err)
		return nil, errors.New("internal error")
	}
	return rule, nil
}

func (s SpecialRouteService) Delete(routeType string, id int) (*domain.SpecialRouteRule, error) {
	normalizedType, err := domain.NormalizeSpecialRouteType(routeType)
	if err != nil {
		return nil, err
	}

	rule, err := s.Repo.FindSpecialRouteRuleByID(id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	if rule == nil || rule.RouteType != normalizedType {
		return nil, errors.New("special route not found")
	}

	if err := s.Repo.DeleteSpecialRouteRule(id); err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	return rule, nil
}
