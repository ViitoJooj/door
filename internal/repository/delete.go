package repository

func (r *SQLite) DeleteUserByID(id int) error {
	_, err := r.db.Exec(`
		DELETE FROM users
		WHERE id = ?
	`, id)

	return err
}

func (r *SQLite) DeleteApplicationByID(id int) error {
	_, err := r.db.Exec(`
		DELETE FROM applications
		WHERE id = ?
	`, id)

	return err
}

func (r *SQLite) DeleteCors(id int) error {
	_, err := r.db.Exec(`DELETE FROM cors WHERE id = ?`, id)
	return err
}

func (r *SQLite) DeleteWhitelistedIP(id int) error {
	_, err := r.db.Exec(`DELETE FROM ip_whitelist WHERE id = ?`, id)
	return err
}

func (r *SQLite) DeleteBlacklistedIP(id int) error {
	_, err := r.db.Exec(`DELETE FROM ip_blacklist WHERE id = ?`, id)
	return err
}

func (r *SQLite) DeleteSpecialRouteRule(id int) error {
	_, err := r.db.Exec(`DELETE FROM special_route_rules WHERE id = ?`, id)
	return err
}

func (r *SQLite) DeleteRouteRule(id int) error {
	_, err := r.db.Exec(`DELETE FROM route_rules WHERE id = ?`, id)
	return err
}
