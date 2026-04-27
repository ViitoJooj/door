CREATE TABLE IF NOT EXISTS rate_limit_settings_backup (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    requests_per_second REAL NOT NULL CHECK (requests_per_second > 0),
    burst INTEGER NOT NULL CHECK (burst > 0),
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO rate_limit_settings_backup (id, requests_per_second, burst, updated_at, created_at)
SELECT id, requests_per_second, burst, updated_at, created_at
FROM rate_limit_settings;

DROP TABLE rate_limit_settings;

ALTER TABLE rate_limit_settings_backup RENAME TO rate_limit_settings;
