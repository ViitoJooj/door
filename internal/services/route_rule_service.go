package services

import (
	"errors"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/repository"
)

type RouteRuleService struct {
	repo repository.RouteRuleRepository
}

func NewRouteRuleService(repo repository.RouteRuleRepository) *RouteRuleService {
	return &RouteRuleService{repo: repo}
}

func (s *RouteRuleService) List() ([]*domain.RouteRule, error) {
	return s.repo.ListRouteRules()
}

func (s *RouteRuleService) GetByID(id int) (*domain.RouteRule, error) {
	rule, err := s.repo.FindRouteRuleByID(id)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, errors.New("route rule not found")
	}
	return rule, nil
}

func (s *RouteRuleService) Create(path, method string, rateLimitEnabled bool, rateLimitRPS float64, rateLimitBurst int, targetURL string, geoRoutingEnabled bool, enabled bool, userID int) (*domain.RouteRule, error) {
	rule, err := domain.NewRouteRule(path, method, rateLimitEnabled, rateLimitRPS, rateLimitBurst, targetURL, geoRoutingEnabled, enabled)
	if err != nil {
		return nil, err
	}
	rule.CreatedBy = userID
	rule.UpdatedBy = userID

	if err := s.repo.CreateRouteRule(rule); err != nil {
		return nil, errors.New("internal error")
	}
	return rule, nil
}

func (s *RouteRuleService) Update(id int, path, method string, rateLimitEnabled bool, rateLimitRPS float64, rateLimitBurst int, targetURL string, geoRoutingEnabled bool, enabled bool, userID int) (*domain.RouteRule, error) {
	existing, err := s.repo.FindRouteRuleByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("route rule not found")
	}

	rule, err := domain.NewRouteRule(path, method, rateLimitEnabled, rateLimitRPS, rateLimitBurst, targetURL, geoRoutingEnabled, enabled)
	if err != nil {
		return nil, err
	}
	rule.ID = id
	rule.CreatedBy = existing.CreatedBy
	rule.UpdatedBy = userID

	if err := s.repo.UpdateRouteRule(rule); err != nil {
		return nil, errors.New("internal error")
	}
	return rule, nil
}

func (s *RouteRuleService) Delete(id int) (*domain.RouteRule, error) {
	rule, err := s.repo.FindRouteRuleByID(id)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, errors.New("route rule not found")
	}
	if err := s.repo.DeleteRouteRule(id); err != nil {
		return nil, errors.New("internal error")
	}
	return rule, nil
}
