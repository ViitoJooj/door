package repository

import "github.com/ViitoJooj/door/internal/domain"

func (r *SQLiteUserRepository) ListUsers() ([]*domain.User, error) {
	rows, err := r.db.Query(`
		SELECT id, username, email, password, updated_at, created_at
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		user := &domain.User{}

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Updated_at,
			&user.Created_at,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *SQLiteUserRepository) FindUserByUsername(username string) (*domain.User, error) {
	row := r.db.QueryRow(`
		SELECT id, username, email, password, updated_at, created_at
		FROM users
		WHERE username = ?
	`, username)

	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Updated_at,
		&user.Created_at,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *SQLiteUserRepository) FindUserByEmail(email string) (*domain.User, error) {
	row := r.db.QueryRow(`
		SELECT id, username, email, password, updated_at, created_at
		FROM users
		WHERE email = ?
	`, email)

	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Updated_at,
		&user.Created_at,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *SQLiteUserRepository) FindUserByID(id float64) (*domain.User, error) {
	row := r.db.QueryRow(`
		SELECT id, username, email, password, updated_at, created_at
		FROM users
		WHERE id = ?
	`, id)

	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Updated_at,
		&user.Created_at,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
