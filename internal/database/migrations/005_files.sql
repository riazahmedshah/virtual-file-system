-- Write your migrate up statements here

CREATE TABLE files (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  name TEXT NOT NULL,
  user_id UUID NOT NULL,
  dir_id UUID NOT NULL,
  gcs_key TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

  CONSTRAINT fk_file_user
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE,

  CONSTRAINT fk_file_dir
    FOREIGN KEY (dir_id)
    REFERENCES dirs (id)
    ON DELETE CASCADE
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
