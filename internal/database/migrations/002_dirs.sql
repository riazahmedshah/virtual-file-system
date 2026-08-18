-- Write your migrate up statements here

CREATE TABLE dirs (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  name TEXT NOT NULL,
  user_id UUID NOT NULL,
  parent_id UUID NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

  CONSTRAINT fk_dir_user
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE,
  
  CONSTRAINT fk_dir_dir
    FOREIGN KEY (parent_id)
    REFERENCES dirs (id)
    ON DELETE CASCADE
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
