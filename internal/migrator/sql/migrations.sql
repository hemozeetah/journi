-- version: 0.01
-- description: create users table

CREATE TABLE IF NOT EXISTS users (
  user_id UUID NOT NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  role TEXT NOT NULL,
  profile TEXT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,

  PRIMARY KEY (user_id),
  UNIQUE (email)
);

-- version: 0.02
-- description: create cities table

CREATE TABLE IF NOT EXISTS cities (
  city_id UUID NOT NULL,
  name TEXT NOT NULL,
  caption TEXT NOT NULL,
  images TEXT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,

  PRIMARY KEY (city_id)
);
