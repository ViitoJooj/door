package repository

import "github.com/ViitoJooj/door/internal/domain"

func (r *SQLiteUserRepository) CreateUser(user *domain.User) error {
	_, err := r.db.Exec(`INSERT INTO users (id, username, email, password, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.Updated_at,
		user.Created_at,
	)
	return err
}
