package repository

import "github.com/ViitoJooj/door/internal/domain"

func (r *SQLiteUserRepository) UpdateUser(user *domain.User) error {
	_, err := r.db.Exec(`
		UPDATE users
		SET username = ?, email = ?, password = ?, updated_at = ?
		WHERE id = ?
	`,
		user.Username,
		user.Email,
		user.Password,
		user.Updated_at,
		user.ID,
	)

	return err
}
