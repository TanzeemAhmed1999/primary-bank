CREATE TABLE users (
  username varchar PRIMARY KEY,
  password varchar NOT NULL,
  full_name varchar NOT NULL,
  email varchar UNIQUE NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now()),
  updated_at timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE accounts
ADD CONSTRAINT fk_accounts_owner
FOREIGN KEY (owner) REFERENCES users (username)
ON DELETE CASCADE;

CREATE UNIQUE INDEX idx_accounts_owner_currency ON accounts (owner, currency);