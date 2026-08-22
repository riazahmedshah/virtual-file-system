-- Write your migrate up statements here

ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

ALTER TABLE users
  ADD COLUMN is_guest BOOLEAN NOT NULL DEFAULT false,
  ADD COLUMN max_storage_limit BIGINT NOT NULL DEFAULT 524288000,
  ADD COLUMN max_file_limit BIGINT NOT NULL DEFAULT 52428800;

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
