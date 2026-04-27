CREATE TABLE IF NOT EXISTS protocol_settings (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    allowed_protocol TEXT NOT NULL CHECK (allowed_protocol IN ('http', 'https', 'both')),
    apply_scope TEXT NOT NULL CHECK (apply_scope IN ('all', 'external')),
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO protocol_settings (id, allowed_protocol, apply_scope)
VALUES (1, 'both', 'all');
