CREATE INDEX IF NOT EXISTS movies_title_idx ON comics USING GIN (to_tsvector('simple', title));