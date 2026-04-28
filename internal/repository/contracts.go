package repository

import (
	"database/sql"
	"time"

	"github.com/ViitoJooj/ward/internal/domain"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) (DotEnvRepository, UserRepository, ApplicationRepository, RequestLogRepository, CorsRepository, RateLimitRepository, IPAccessListRepository, ProtocolSettingsRepository, SpecialRouteRepository) {
	repo := &SQLite{db: db}
	return repo, repo, repo, repo, repo, repo, repo, repo, repo
}

type UserRepository interface {
	CreateUser(user *domain.User) error
	CountUsers() (int, error)
	FindUserByID(id int) (*domain.User, error)
	FindUserByEmail(email string) (*domain.User, error)
	FindUserByUsername(username string) (*domain.User, error)
	ListUsers() ([]*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUserByID(id int) error
}

type ApplicationRepository interface {
	CreateApplication(application *domain.Application) error
	FindApplicationByID(id int) (*domain.Application, error)
	FindApplicationByURL(url string) (*domain.Application, error)
	FindApplicationByCountry(country string) (*domain.Application, error)
	GetRandomApplication() (*domain.Application, error)
	ListApplications() ([]*domain.Application, error)
	UpdateApplication(application *domain.Application) error
	DeleteApplicationByID(id int) error
}

type RequestLogRepository interface {
	InsertRequestLog(log *domain.RequestLog) error
	ListRequestLogs() ([]*domain.RequestLog, error)
	ListRequestLogsSince(since time.Time, limit int) ([]*domain.RequestLog, error)
}

type DotEnvRepository interface {
	FindVar(id int) (*domain.Env, error)
	ChangeVar(*domain.Env) error
	GetAllVars() ([]*domain.Env, error)
}

type CorsRepository interface {
	FindAllCors() ([]*domain.Cors, error)
	FindCorsByID(id int) (*domain.Cors, error)
	CreateCors(*domain.Cors) error
	ChangeCors(*domain.Cors) error
	DeleteCors(id int) error
}

type RateLimitRepository interface {
	GetRateLimitSettings() (*domain.RateLimitSettings, error)
	UpsertRateLimitSettings(*domain.RateLimitSettings) error
}

type IPAccessListRepository interface {
	ListWhitelistedIPs() ([]*domain.IPAccessEntry, error)
	FindWhitelistedIPByID(id int) (*domain.IPAccessEntry, error)
	CreateWhitelistedIP(*domain.IPAccessEntry) error
	UpdateWhitelistedIP(*domain.IPAccessEntry) error
	DeleteWhitelistedIP(id int) error

	ListBlacklistedIPs() ([]*domain.IPAccessEntry, error)
	FindBlacklistedIPByID(id int) (*domain.IPAccessEntry, error)
	CreateBlacklistedIP(*domain.IPAccessEntry) error
	UpdateBlacklistedIP(*domain.IPAccessEntry) error
	DeleteBlacklistedIP(id int) error
}

type ProtocolSettingsRepository interface {
	GetProtocolSettings() (*domain.ProtocolSettings, error)
	UpsertProtocolSettings(*domain.ProtocolSettings) error
}

type SpecialRouteRepository interface {
	ListSpecialRouteRules(routeType string) ([]*domain.SpecialRouteRule, error)
	FindSpecialRouteRuleByID(id int) (*domain.SpecialRouteRule, error)
	CreateSpecialRouteRule(*domain.SpecialRouteRule) error
	UpdateSpecialRouteRule(*domain.SpecialRouteRule) error
	DeleteSpecialRouteRule(id int) error
}
