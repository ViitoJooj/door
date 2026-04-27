CREATE TABLE IF NOT EXISTS ip_whitelist_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip TEXT NOT NULL UNIQUE,
    created_by INTEGER NOT NULL DEFAULT 0,
    updated_by INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO ip_whitelist_new (id, ip, created_by, updated_by, created_at, updated_at)
SELECT id, ip, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
FROM ip_whitelist;

DROP TABLE ip_whitelist;
ALTER TABLE ip_whitelist_new RENAME TO ip_whitelist;

CREATE TABLE IF NOT EXISTS ip_blacklist_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip TEXT NOT NULL UNIQUE,
    created_by INTEGER NOT NULL DEFAULT 0,
    updated_by INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO ip_blacklist_new (id, ip, created_by, updated_by, created_at, updated_at)
SELECT id, ip, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
FROM ip_blacklist;

DROP TABLE ip_blacklist;
ALTER TABLE ip_blacklist_new RENAME TO ip_blacklist;
