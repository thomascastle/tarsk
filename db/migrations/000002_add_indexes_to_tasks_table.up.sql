CREATE INDEX IF NOT EXISTS tasks_description_idx ON tasks USING GIN (to_tsvector('simple', description));
