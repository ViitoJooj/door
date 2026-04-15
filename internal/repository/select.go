package repository

import (
	"database/sql"

	"github.com/ViitoJooj/door/internal/domain"
)

func (r *SQLite) ListUsers() ([]*domain.User, error) {
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

func (r *SQLite) FindUserByUsername(username string) (*domain.User, error) {
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

func (r *SQLite) FindUserByEmail(email string) (*domain.User, error) {
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

func (r *SQLite) FindUserByID(id int) (*domain.User, error) {
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

func (r *SQLite) ListApplications() ([]*domain.Application, error) {
	rows, err := r.db.Query(`
		SELECT id, url, country, created_by, updated_at, created_at
		FROM applications
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []*domain.Application

	for rows.Next() {
		application := &domain.Application{}

		err := rows.Scan(
			&application.ID,
			&application.Url,
			&application.Country,
			&application.Created_by,
			&application.Updated_at,
			&application.Created_at,
		)
		if err != nil {
			return nil, err
		}

		applications = append(applications, application)
	}

	return applications, rows.Err()
}

func (r *SQLite) FindApplicationByID(id int) (*domain.Application, error) {
	application := &domain.Application{}

	err := r.db.QueryRow(`
		SELECT id, url, country, created_by, updated_at, created_at
		FROM applications
		WHERE id = ?
	`, id).Scan(
		&application.ID,
		&application.Url,
		&application.Country,
		&application.Created_by,
		&application.Updated_at,
		&application.Created_at,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return application, err
}

func (r *SQLite) FindApplicationByCountry(country string) (*domain.Application, error) {
	application := &domain.Application{}

	err := r.db.QueryRow(`
		SELECT id, url, country, created_by, updated_at, created_at
		FROM applications
		WHERE country = ?
	`, country).Scan(
		&application.ID,
		&application.Url,
		&application.Country,
		&application.Created_by,
		&application.Updated_at,
		&application.Created_at,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return application, err
}

func (r *SQLite) ListRequestLogs() ([]*domain.RequestLog, error) {
	rows, err := r.db.Query(`
		SELECT id, method, path, query_string, status_code, response_time_ms, ip, country, user_agent, referer, request_size, response_size, internal, created_at
		FROM request_logs
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.RequestLog

	for rows.Next() {
		entry := &domain.RequestLog{}

		err := rows.Scan(
			&entry.ID,
			&entry.Method,
			&entry.Path,
			&entry.QueryString,
			&entry.StatusCode,
			&entry.ResponseTimeMs,
			&entry.IP,
			&entry.Country,
			&entry.UserAgent,
			&entry.Referer,
			&entry.RequestSize,
			&entry.ResponseSize,
			&entry.Internal,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		logs = append(logs, entry)
	}

	return logs, rows.Err()
}

func (r *SQLite) FindApplicationByURL(url string) (*domain.Application, error) {
	application := &domain.Application{}

	err := r.db.QueryRow(`
		SELECT id, url, country, created_by, updated_at, created_at
		FROM applications
		WHERE url = ?
	`, url).Scan(
		&application.ID,
		&application.Url,
		&application.Country,
		&application.Created_by,
		&application.Updated_at,
		&application.Created_at,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return application, err
}
