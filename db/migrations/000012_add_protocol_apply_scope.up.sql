ALTER TABLE protocol_settings
ADD COLUMN apply_scope TEXT NOT NULL DEFAULT 'all' CHECK (apply_scope IN ('all', 'external'));
