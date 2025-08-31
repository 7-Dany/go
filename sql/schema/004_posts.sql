-- +goose up
CREATE TABLE posts(
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  url TEXT NOT NULL UNIQUE, 
  feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  published_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);


-- +goose down
DROP TABLE posts;