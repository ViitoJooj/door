CREATE TABLE IF NOT EXISTS Applications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    Url TEXT NOT NULL,
    Country TEXT NOT NULL,
    Created_by INTEGER,
    Updated_at DATETIME,
    Created_at DATETIME
);