CREATE TABLE IF NOT EXISTS special_route_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    route_type TEXT NOT NULL CHECK (route_type IN ('login', 'register')),
    path TEXT NOT NULL,
    max_distinct_requests INTEGER NOT NULL CHECK (max_distinct_requests > 0),
    window_seconds INTEGER NOT NULL CHECK (window_seconds > 0),
    ban_seconds INTEGER NOT NULL CHECK (ban_seconds > 0),
    enabled INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
    created_by INTEGER NOT NULL,
    updated_by INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(route_type, path),
    FOREIGN KEY(created_by) REFERENCES users(id),
    FOREIGN KEY(updated_by) REFERENCES users(id)
);
