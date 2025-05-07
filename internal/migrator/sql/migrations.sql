-- version: 0.01
-- description: create table users

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
