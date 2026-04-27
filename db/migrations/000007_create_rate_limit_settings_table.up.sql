CREATE TABLE IF NOT EXISTS rate_limit_settings (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    requests_per_second REAL NOT NULL CHECK (requests_per_second > 0),
    burst INTEGER NOT NULL CHECK (burst > 0),
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO rate_limit_settings (id, requests_per_second, burst)
VALUES (1, 1, 5);
