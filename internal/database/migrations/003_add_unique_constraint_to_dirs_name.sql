-- Write your migrate up statements here

ALTER TABLE dirs
ADD CONSTRAINT unique_user_parent_dir_name UNIQUE (user_id, parent_id, name)

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
