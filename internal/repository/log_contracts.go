package repository

import (
	"database/sql"

	"github.com/ViitoJooj/door/internal/domain"
)

type RequestLogRepository interface {
	InsertRequestLog(log *domain.RequestLog) error
}

type SQLiteRequestLogRepository struct {
	db *sql.DB
}

func NewSQLiteRequestLogRepository(db *sql.DB) RequestLogRepository {
	return &SQLiteRequestLogRepository{db: db}
}

func (r *SQLiteRequestLogRepository) InsertRequestLog(log *domain.RequestLog) error {
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
