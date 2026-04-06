package repository

func (r *SQLiteUserRepository) DeleteUserByID(id int) error {
	_, err := r.db.Exec(`
		DELETE FROM users
		WHERE id = ?
	`, id)

	return err
}
