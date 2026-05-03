CREATE TABLE IF NOT EXISTS route_rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL,
    method TEXT NOT NULL DEFAULT '',
    rate_limit_enabled INTEGER NOT NULL DEFAULT 0 CHECK (rate_limit_enabled IN (0, 1)),
    rate_limit_rps REAL NOT NULL DEFAULT 0,
    rate_limit_burst INTEGER NOT NULL DEFAULT 0,
    target_url TEXT NOT NULL DEFAULT '',
    geo_routing_enabled INTEGER NOT NULL DEFAULT 0 CHECK (geo_routing_enabled IN (0, 1)),
    enabled INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
    created_by INTEGER NOT NULL,
    updated_by INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(path, method),
    FOREIGN KEY(created_by) REFERENCES users(id),
    FOREIGN KEY(updated_by) REFERENCES users(id)
);
