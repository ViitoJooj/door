package services

import (
	"errors"
	"log"
	"net"
	"strings"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/internal/repository"
)

type IPAccessListService struct {
	Repo     repository.IPAccessListRepository
	AuthRepo repository.UserRepository
}

func NewIPAccessListService(repo repository.IPAccessListRepository, authRepo repository.UserRepository) *IPAccessListService {
	return &IPAccessListService{
		Repo:     repo,
		AuthRepo: authRepo,
	}
}

func normalizeIP(ip string) string {
	return strings.TrimSpace(ip)
}

func validateIP(ip string) error {
	if ip == "" {
		return errors.New("ip is required")
	}
	if parsed := net.ParseIP(ip); parsed == nil {
		return errors.New("invalid ip")
	}
	return nil
}

func (s IPAccessListService) GetWhitelist() ([]*domain.IPAccessEntry, error) {
	return s.Repo.ListWhitelistedIPs()
}

func (s IPAccessListService) GetBlacklist() ([]*domain.IPAccessEntry, error) {
	return s.Repo.ListBlacklistedIPs()
}

func (s IPAccessListService) ensureUser(userID int) error {
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

func (s IPAccessListService) CreateWhitelist(ip string, userID int) (*domain.IPAccessEntry, error) {
	ip = normalizeIP(ip)
	if err := validateIP(ip); err != nil {
		return nil, err
	}
	if err := s.ensureUser(userID); err != nil {
		return nil, err
	}

	entry := &domain.IPAccessEntry{IP: ip, CreatedBy: userID, UpdatedBy: userID}
	if err := s.Repo.CreateWhitelistedIP(entry); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, errors.New("ip already exists")
		}
		log.Println(err)
		return nil, errors.New("internal error")
	}
	return entry, nil
}

func (s IPAccessListService) CreateBlacklist(ip string, userID int) (*domain.IPAccessEntry, error) {
	ip = normalizeIP(ip)
	if err := validateIP(ip); err != nil {
		return nil, err
	}
	if err := s.ensureUser(userID); err != nil {
		return nil, err
	}

	entry := &domain.IPAccessEntry{IP: ip, CreatedBy: userID, UpdatedBy: userID}
	if err := s.Repo.CreateBlacklistedIP(entry); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, errors.New("ip already exists")
		}
		log.Println(err)
		return nil, errors.New("internal error")
	}
	return entry, nil
}

func (s IPAccessListService) UpdateWhitelist(id int, ip string, userID int) (*domain.IPAccessEntry, error) {
	ip = normalizeIP(ip)
	if err := validateIP(ip); err != nil {
		return nil, err
	}
	if err := s.ensureUser(userID); err != nil {
		return nil, err
	}

	entry, err := s.Repo.FindWhitelistedIPByID(id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	if entry == nil {
		return nil, errors.New("ip not found")
	}

	entry.IP = ip
	entry.UpdatedBy = userID
	if err := s.Repo.UpdateWhitelistedIP(entry); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, errors.New("ip already exists")
		}
		log.Println(err)
		return nil, errors.New("internal error")
	}
	return entry, nil
}

func (s IPAccessListService) UpdateBlacklist(id int, ip string, userID int) (*domain.IPAccessEntry, error) {
	ip = normalizeIP(ip)
	if err := validateIP(ip); err != nil {
		return nil, err
	}
	if err := s.ensureUser(userID); err != nil {
		return nil, err
	}

	entry, err := s.Repo.FindBlacklistedIPByID(id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	if entry == nil {
		return nil, errors.New("ip not found")
	}

	entry.IP = ip
	entry.UpdatedBy = userID
	if err := s.Repo.UpdateBlacklistedIP(entry); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, errors.New("ip already exists")
		}
		log.Println(err)
		return nil, errors.New("internal error")
	}
	return entry, nil
}

func (s IPAccessListService) DeleteWhitelist(id int) (*domain.IPAccessEntry, error) {
	entry, err := s.Repo.FindWhitelistedIPByID(id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	if entry == nil {
		return nil, errors.New("ip not found")
	}

	if err := s.Repo.DeleteWhitelistedIP(id); err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	return entry, nil
}

func (s IPAccessListService) DeleteBlacklist(id int) (*domain.IPAccessEntry, error) {
	entry, err := s.Repo.FindBlacklistedIPByID(id)
	if err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	if entry == nil {
		return nil, errors.New("ip not found")
	}

	if err := s.Repo.DeleteBlacklistedIP(id); err != nil {
		log.Println(err)
		return nil, errors.New("internal error")
	}
	return entry, nil
}
