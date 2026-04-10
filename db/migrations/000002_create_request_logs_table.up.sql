CREATE TABLE IF NOT EXISTS request_logs (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      method TEXT NOT NULL,
      path TEXT NOT NULL,
      query_string TEXT,
      status_code INTEGER NOT NULL,
      response_time_ms INTEGER NOT NULL,
      ip TEXT,
      country TEXT,
      user_agent TEXT,
      referer TEXT,
      request_size INTEGER,
      response_size INTEGER,
      internal BOOL,
      created_at DATETIME DEFAULT CURRENT_TIMESTAMP
  );
