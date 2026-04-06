package domain

import (
	"errors"
	"time"
)

type User struct {
	ID         int
	Username   string
	Email      string
	Password   string
	Updated_at time.Time
	Created_at time.Time
}

func NewUser(username string, email string, password string) (*User, error) {
	// Username validator
	if len(username) > 250 {
		return nil, errors.New("Name is too large.")
	} else if len(username) < 3 {
		return nil, errors.New("username is too short")
	}

	// Email validator
	if len(email) > 250 {
		return nil, errors.New("Email is too large.")
	} else if len(email) < 10 {
		return nil, errors.New("Email is too short.")
	}

	// Password validator
	if len(password) > 50 {
		return nil, errors.New("Password is too large.")
	} else if len(password) < 6 {
		return nil, errors.New("Password is too short.")
	}

	user := &User{
		Username:   username,
		Email:      email,
		Password:   password,
		Updated_at: time.Now(),
		Created_at: time.Now(),
	}

	return user, nil
}
