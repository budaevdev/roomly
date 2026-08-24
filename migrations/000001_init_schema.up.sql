CREATE EXTENSION IF NOT EXISTS btree_gist;

 CREATE TABLE users (
  id serial PRIMARY KEY,
  username text UNIQUE NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE listings (
  id serial PRIMARY KEY,
  title text NOT NULL,
  description text NOT NULL,
  owner_id integer NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  FOREIGN KEY (owner_id) REFERENCES users (id)
);

CREATE TABLE bookings (
  id serial PRIMARY KEY,
  listing_id integer NOT NULL,
  guest_id integer NOT NULL,
  during daterange NOT NULL,
  status text NOT NULL DEFAULT 'active',
  EXCLUDE USING gist (listing_id WITH =, during WITH &&) WHERE (status = 'active'),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  FOREIGN KEY (listing_id) REFERENCES listings (id),
  FOREIGN KEY (guest_id) REFERENCES users (id)
);