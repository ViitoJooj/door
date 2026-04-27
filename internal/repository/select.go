package repository

import (
	"database/sql"

	"github.com/ViitoJooj/ward/internal/domain"
	"github.com/ViitoJooj/ward/pkg/initializer"
)

func (r *SQLite) FindVar(id int) (*domain.Env, error) {
	row := r.db.QueryRow(`SELECT id, name, value FROM env WHERE id = ?`, id)
	env := &domain.Env{}
	err := row.Scan(&env.Id, &env.Name, &env.Value)
	if err != nil {
		return nil, err
	}

	if !initializer.IsMasterKeyVar(env.Name) {
		value, err := initializer.DecryptValue(env.Value)
		if err != nil {
			return nil, err
		}
		env.Value = value
	}

	return env, nil
}

func (r *SQLite) GetAllVars() ([]*domain.Env, error) {
	rows, err := r.db.Query(`
		SELECT id, name, value
		FROM env
		ORDER BY CASE WHEN name = 'MASTER_KEY' THEN 0 ELSE 1 END, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envs []*domain.Env

	for rows.Next() {
		env := &domain.Env{}
		err := rows.Scan(
			&env.Id,
			&env.Name,
			&env.Value,
		)
		if err != nil {
			return nil, err
		}

		if !initializer.IsMasterKeyVar(env.Name) {
			value, decErr := initializer.DecryptValue(env.Value)
			if decErr != nil {
				return nil, decErr
			}
			env.Value = value
		}

		envs = append(envs, env)

	}

	return envs, nil
}

func (r *SQLite) ListUsers() ([]*domain.User, error) {
	rows, err := r.db.Query(`
		SELECT id, username, email, password, role, active, updated_at, created_at
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
			&user.Role,
			&user.Active,
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

func (r *SQLite) CountUsers() (int, error) {
	var total int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *SQLite) FindUserByUsername(username string) (*domain.User, error) {
	row := r.db.QueryRow(`
		SELECT id, username, email, password, role, active, updated_at, created_at
		FROM users
		WHERE username = ?
	`, username)

	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Active,
		&user.Updated_at,
		&user.Created_at,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *SQLite) FindUserByEmail(email string) (*domain.User, error) {
	row := r.db.QueryRow(`
		SELECT id, username, email, password, role, active, updated_at, created_at
		FROM users
		WHERE email = ?
	`, email)

	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Active,
		&user.Updated_at,
		&user.Created_at,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *SQLite) FindUserByID(id int) (*domain.User, error) {
	row := r.db.QueryRow(`
		SELECT id, username, email, password, role, active, updated_at, created_at
		FROM users
		WHERE id = ?
	`, id)

	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Active,
		&user.Updated_at,
		&user.Created_at,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
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

func (r *SQLite) GetRandomApplication() (*domain.Application, error) {
	row := r.db.QueryRow(`
		SELECT id, url, country, created_by, updated_at, created_at
		FROM applications
		ORDER BY RANDOM()
		LIMIT 1
	`)

	application := &domain.Application{}

	err := row.Scan(
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

	return application, nil
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

func (r *SQLite) FindAllCors() ([]*domain.Cors, error) {
	rows, err := r.db.Query(`SELECT id, origin FROM cors`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var corsList []*domain.Cors

	for rows.Next() {
		var cors domain.Cors
		if err := rows.Scan(&cors.Id, &cors.Origin); err != nil {
			return nil, err
		}
		corsList = append(corsList, &cors)
	}

	return corsList, nil
}

func (r *SQLite) FindCorsByID(id int) (*domain.Cors, error) {
	var cors domain.Cors

	err := r.db.QueryRow(`SELECT id, origin FROM cors WHERE id = ?`, id).
		Scan(&cors.Id, &cors.Origin)

	if err != nil {
		return nil, err
	}

	return &cors, nil
}

func (r *SQLite) GetRateLimitSettings() (*domain.RateLimitSettings, error) {
	settings := &domain.RateLimitSettings{}
	var progressiveEnabled int

	err := r.db.QueryRow(`
		SELECT id, requests_per_second, burst, progressive_enabled, updated_at, created_at
		FROM rate_limit_settings
		WHERE id = 1
	`).Scan(
		&settings.ID,
		&settings.RequestsPerSecond,
		&settings.Burst,
		&progressiveEnabled,
		&settings.UpdatedAt,
		&settings.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	settings.Progressive = progressiveEnabled == 1

	return settings, nil
}

func (r *SQLite) ListWhitelistedIPs() ([]*domain.IPAccessEntry, error) {
	rows, err := r.db.Query(`
		SELECT id, ip, created_by, updated_by, created_at, updated_at
		FROM ip_whitelist
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]*domain.IPAccessEntry, 0)
	for rows.Next() {
		entry := &domain.IPAccessEntry{}
		if err := rows.Scan(
			&entry.ID,
			&entry.IP,
			&entry.CreatedBy,
			&entry.UpdatedBy,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (r *SQLite) FindWhitelistedIPByID(id int) (*domain.IPAccessEntry, error) {
	entry := &domain.IPAccessEntry{}
	err := r.db.QueryRow(`
		SELECT id, ip, created_by, updated_by, created_at, updated_at
		FROM ip_whitelist
		WHERE id = ?
	`, id).Scan(
		&entry.ID,
		&entry.IP,
		&entry.CreatedBy,
		&entry.UpdatedBy,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (r *SQLite) ListBlacklistedIPs() ([]*domain.IPAccessEntry, error) {
	rows, err := r.db.Query(`
		SELECT id, ip, created_by, updated_by, created_at, updated_at
		FROM ip_blacklist
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]*domain.IPAccessEntry, 0)
	for rows.Next() {
		entry := &domain.IPAccessEntry{}
		if err := rows.Scan(
			&entry.ID,
			&entry.IP,
			&entry.CreatedBy,
			&entry.UpdatedBy,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (r *SQLite) FindBlacklistedIPByID(id int) (*domain.IPAccessEntry, error) {
	entry := &domain.IPAccessEntry{}
	err := r.db.QueryRow(`
		SELECT id, ip, created_by, updated_by, created_at, updated_at
		FROM ip_blacklist
		WHERE id = ?
	`, id).Scan(
		&entry.ID,
		&entry.IP,
		&entry.CreatedBy,
		&entry.UpdatedBy,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (r *SQLite) GetProtocolSettings() (*domain.ProtocolSettings, error) {
	settings := &domain.ProtocolSettings{}

	err := r.db.QueryRow(`
		SELECT id, allowed_protocol, apply_scope, updated_at, created_at
		FROM protocol_settings
		WHERE id = 1
	`).Scan(
		&settings.ID,
		&settings.AllowedProtocol,
		&settings.ApplyScope,
		&settings.UpdatedAt,
		&settings.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (r *SQLite) ListSpecialRouteRules(routeType string) ([]*domain.SpecialRouteRule, error) {
	rows, err := r.db.Query(`
		SELECT id, route_type, path, max_distinct_requests, window_seconds, ban_seconds, enabled, created_by, updated_by, created_at, updated_at
		FROM special_route_rules
		WHERE route_type = ?
		ORDER BY id DESC
	`, routeType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := make([]*domain.SpecialRouteRule, 0)
	for rows.Next() {
		rule := &domain.SpecialRouteRule{}
		var enabledInt int
		if err := rows.Scan(
			&rule.ID,
			&rule.RouteType,
			&rule.Path,
			&rule.MaxDistinctRequests,
			&rule.WindowSeconds,
			&rule.BanSeconds,
			&enabledInt,
			&rule.CreatedBy,
			&rule.UpdatedBy,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		); err != nil {
			return nil, err
		}
		rule.Enabled = enabledInt == 1
		rules = append(rules, rule)
	}

	return rules, rows.Err()
}

func (r *SQLite) FindSpecialRouteRuleByID(id int) (*domain.SpecialRouteRule, error) {
	rule := &domain.SpecialRouteRule{}
	var enabledInt int
	err := r.db.QueryRow(`
		SELECT id, route_type, path, max_distinct_requests, window_seconds, ban_seconds, enabled, created_by, updated_by, created_at, updated_at
		FROM special_route_rules
		WHERE id = ?
	`, id).Scan(
		&rule.ID,
		&rule.RouteType,
		&rule.Path,
		&rule.MaxDistinctRequests,
		&rule.WindowSeconds,
		&rule.BanSeconds,
		&enabledInt,
		&rule.CreatedBy,
		&rule.UpdatedBy,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rule.Enabled = enabledInt == 1
	return rule, nil
}
