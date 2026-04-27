package repository

import (
	"database/sql"

	"github.com/ViitoJooj/ward/internal/domain"
)

func (r *SQLite) CreateUser(user *domain.User) error {
	_, err := r.db.Exec(`INSERT INTO users (username, email, password, role, active, updated_at, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.Username,
		user.Email,
		user.Password,
		user.Role,
		user.Active,
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

func (r *SQLite) CreateCors(cors *domain.Cors) error {
	_, err := r.db.Exec(`INSERT INTO cors (name, origin) VALUES (?, ?)`,
		cors.Name, cors.Origin)
	return err
}

func (r *SQLite) CreateWhitelistedIP(entry *domain.IPAccessEntry) error {
	res, err := r.db.Exec(`
		INSERT INTO ip_whitelist (ip, created_by, updated_by)
		VALUES (?, ?, ?)
	`, entry.IP, entry.CreatedBy, entry.UpdatedBy)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	entry.ID = int(id)
	return r.db.QueryRow(`
		SELECT id, ip, created_by, updated_by, created_at, updated_at
		FROM ip_whitelist
		WHERE id = ?
	`, entry.ID).Scan(
		&entry.ID,
		&entry.IP,
		&entry.CreatedBy,
		&entry.UpdatedBy,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)
}

func (r *SQLite) CreateBlacklistedIP(entry *domain.IPAccessEntry) error {
	res, err := r.db.Exec(`
		INSERT INTO ip_blacklist (ip, created_by, updated_by)
		VALUES (?, ?, ?)
	`, entry.IP, entry.CreatedBy, entry.UpdatedBy)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	entry.ID = int(id)
	return r.db.QueryRow(`
		SELECT id, ip, created_by, updated_by, created_at, updated_at
		FROM ip_blacklist
		WHERE id = ?
	`, entry.ID).Scan(
		&entry.ID,
		&entry.IP,
		&entry.CreatedBy,
		&entry.UpdatedBy,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)
}

func (r *SQLite) CreateSpecialRouteRule(rule *domain.SpecialRouteRule) error {
	enabledInt := 0
	if rule.Enabled {
		enabledInt = 1
	}

	res, err := r.db.Exec(`
		INSERT INTO special_route_rules (route_type, path, max_distinct_requests, window_seconds, ban_seconds, enabled, created_by, updated_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		rule.RouteType,
		rule.Path,
		rule.MaxDistinctRequests,
		rule.WindowSeconds,
		rule.BanSeconds,
		enabledInt,
		rule.CreatedBy,
		rule.UpdatedBy,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	createdRule, err := r.FindSpecialRouteRuleByID(int(id))
	if err != nil {
		return err
	}
	if createdRule == nil {
		return sql.ErrNoRows
	}

	*rule = *createdRule
	return nil
}
