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

-- version: 0.03
-- description: create places table

CREATE TABLE IF NOT EXISTS places (
  place_id UUID NOT NULL,
  city_id UUID NOT NULL,
  name TEXT NOT NULL,
  caption TEXT NOT NULL,
  type TEXT NOT NULL,
  images TEXT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,

  PRIMARY KEY (place_id),
	FOREIGN KEY (city_id) REFERENCES cities(city_id) ON DELETE CASCADE
);

-- version: 0.04
-- description: create posts table

CREATE TABLE IF NOT EXISTS posts (
  post_id UUID NOT NULL,
  user_id UUID NOT NULL,
  place_id UUID NOT NULL,
  caption TEXT NOT NULL,
  images TEXT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,

  PRIMARY KEY (post_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
	FOREIGN KEY (place_id) REFERENCES places(place_id) ON DELETE CASCADE
);
