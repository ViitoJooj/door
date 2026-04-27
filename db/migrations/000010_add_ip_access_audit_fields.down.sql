CREATE TABLE IF NOT EXISTS ip_whitelist_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip TEXT NOT NULL UNIQUE
);

INSERT INTO ip_whitelist_old (id, ip)
SELECT id, ip
FROM ip_whitelist;

DROP TABLE ip_whitelist;
ALTER TABLE ip_whitelist_old RENAME TO ip_whitelist;

CREATE TABLE IF NOT EXISTS ip_blacklist_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip TEXT NOT NULL UNIQUE
);

INSERT INTO ip_blacklist_old (id, ip)
SELECT id, ip
FROM ip_blacklist;

DROP TABLE ip_blacklist;
ALTER TABLE ip_blacklist_old RENAME TO ip_blacklist;
