package repository

import "github.com/ViitoJooj/ward/internal/domain"

func (r *SQLite) UpdateUser(user *domain.User) error {
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

func (r *SQLite) UpdateApplication(application *domain.Application) error {
	_, err := r.db.Exec(`
		UPDATE applications
		SET url = ?, country = ?, created_by = ?, updated_at = ?
		WHERE id = ?
	`,
		application.Url,
		application.Country,
		application.Created_by,
		application.Updated_at,
		application.ID,
	)

	return err
}
