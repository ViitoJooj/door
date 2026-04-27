ALTER TABLE rate_limit_settings
ADD COLUMN progressive_enabled INTEGER NOT NULL DEFAULT 0 CHECK (progressive_enabled IN (0, 1));
