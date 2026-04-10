package repository

import "github.com/ViitoJooj/door/internal/domain"

func (r *SQLite) CreateUser(user *domain.User) error {
	_, err := r.db.Exec(`INSERT INTO users (username, email, password, updated_at, created_at) VALUES (?, ?, ?, ?, ?)`,
		user.Username,
		user.Email,
		user.Password,
		user.Updated_at,
		user.Created_at,
	)
	return err
}

func (r *SQLite) CreateApplication(application *domain.Application) error {
	_, err := r.db.Exec(`INSERT INTO applications (url, country, created_by, updated_at, created_at) VALUES (?, ?, ?, ?, ?)`,
		application.Url,
		application.Country,
		application.Created_by,
		application.Updated_at,
		application.Created_at,
	)
	return err
}

func (r *SQLite) InsertRequestLog(log *domain.RequestLog) error {
	_, err := r.db.Exec(`
		INSERT INTO request_logs
			(method, path, query_string, status_code, response_time_ms, ip, country, user_agent, referer, request_size, response_size, internal)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		log.Method,
		log.Path,
		log.QueryString,
		log.StatusCode,
		log.ResponseTimeMs,
		log.IP,
		log.Country,
		log.UserAgent,
		log.Referer,
		log.RequestSize,
		log.ResponseSize,
		log.Internal,
	)
	return err
}
