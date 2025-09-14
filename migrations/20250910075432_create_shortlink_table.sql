-- +goose Up
CREATE TABLE short_links (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  short_id     TEXT NOT NULL UNIQUE,
  original_url TEXT NOT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at   TIMESTAMPTZ NULL,
  max_visits   BIGINT NOT NULL DEFAULT 0,
  visits       BIGINT NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE short_links;