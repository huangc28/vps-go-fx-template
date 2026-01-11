-- Placeholder schema for sqlc.
-- Replace with your real schema (or export it from your DB).

CREATE TABLE IF NOT EXISTS healthchecks (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

