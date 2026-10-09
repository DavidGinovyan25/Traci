CREATE TYPE role AS ENUM (
  'admin',
  'user'
);

CREATE TYPE movie_genre AS ENUM (
  'ACTION',
  'ADVENTURE',
  'ANIMATION',
  'COMEDY',
  'CRIME',
  'DOCUMENTARY',
  'DRAMA',
  'FAMILY',
  'FANTASY',
  'HISTORY',
  'HORROR',
  'MUSIC',
  'MYSTERY',
  'ROMANCE',
  'SCIENCE_FICTION',
  'THRILLER',
  'WAR',
  'WESTERN'
);

CREATE TYPE movie_status AS ENUM (
  'PLANNED',
  'WATCHING',
  'WATCHED',
  'DROPPED'
);

CREATE TABLE users (
  id uuid PRIMARY KEY,
  role role NOT NULL DEFAULT 'user',
  username varchar(50) UNIQUE NOT NULL,
  first_name text NOT NULL,
  second_name text NOT NULL,
  email varchar(255) UNIQUE NOT NULL,
  password_hash varchar(255) NOT NULL,
  is_blocked boolean NOT NULL DEFAULT false
);

CREATE TABLE movies (
  id uuid PRIMARY KEY,
  name text NOT NULL,
  duration_min int NOT NULL CHECK (duration_min > 0),
  genre movie_genre NOT NULL,
  release_year int NOT NULL CHECK (release_year >= 1888),
  description text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE user_movies (
  id uuid PRIMARY KEY,
  movie_id uuid NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  personal_rating int CHECK (personal_rating BETWEEN 1 AND 10),
  review varchar(1000),
  status movie_status NOT NULL DEFAULT 'PLANNED',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),

  UNIQUE (user_id, movie_id)
);
