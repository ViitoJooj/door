package repository

import (
	"database/sql"

	"github.com/ViitoJooj/door/internal/domain"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) (UserRepository, ApplicationRepository, RequestLogRepository) {
	repo := &SQLite{db: db}
	return repo, repo, repo
}

type UserRepository interface {
	CreateUser(user *domain.User) error
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
	ListApplications() ([]*domain.Application, error)
	UpdateApplication(application *domain.Application) error
	DeleteApplicationByID(id int) error
}

type RequestLogRepository interface {
	InsertRequestLog(log *domain.RequestLog) error
}
