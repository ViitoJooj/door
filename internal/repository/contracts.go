package repository

import (
	"database/sql"

	"github.com/ViitoJooj/door/internal/domain"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) UserRepository {
	return &SQLiteUserRepository{db: db}
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
