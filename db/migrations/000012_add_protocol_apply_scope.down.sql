CREATE TABLE IF NOT EXISTS protocol_settings_backup (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    allowed_protocol TEXT NOT NULL CHECK (allowed_protocol IN ('http', 'https', 'both')),
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO protocol_settings_backup (id, allowed_protocol, updated_at, created_at)
SELECT id, allowed_protocol, updated_at, created_at
FROM protocol_settings;

DROP TABLE protocol_settings;
ALTER TABLE protocol_settings_backup RENAME TO protocol_settings;
