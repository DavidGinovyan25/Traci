CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_movies_name_trgm
ON movies USING GIN (name gin_trgm_ops);

CREATE INDEX idx_movies_created_at_id
ON movies (created_at DESC, id ASC);

CREATE INDEX idx_movies_genre_created_at_id
ON movies (genre, created_at DESC, id ASC);

CREATE INDEX idx_user_movies_movie_id
ON user_movies (movie_id);

CREATE INDEX idx_user_movies_user_created_at_id
ON user_movies (user_id, created_at DESC, id ASC);

CREATE INDEX idx_user_movies_user_status_created_at_id
ON user_movies (user_id, status, created_at DESC, id ASC);