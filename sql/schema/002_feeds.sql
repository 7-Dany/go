-- +goose up
CREATE TABLE feeds (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  url TEXT UNIQUE NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
  last_feched_at TIMESTAMP DEFAULT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);


-- +goose down
DROP TABLE feeds;