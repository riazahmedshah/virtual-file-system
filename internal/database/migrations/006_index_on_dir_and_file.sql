-- Write your migrate up statements here

CREATE INDEX idx_dirs_parent_id ON dirs(parent_id);

CREATE INDEX idx_files_dir_id ON files(dir_id);
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
