-- Write your migrate up statements here

ALTER TABLE files
  ADD COLUMN size BIGINT NOT NULL,
  ADD COLUMN ext TEXT NOT NULL

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
