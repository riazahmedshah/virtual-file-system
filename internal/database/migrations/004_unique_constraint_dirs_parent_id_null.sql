-- Write your migrate up statements here

CREATE UNIQUE INDEX idx_unique_root_dir_per_user
ON dirs (user_id)
WHERE parent_id IS NULL

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
