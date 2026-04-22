CREATE TABLE IF NOT EXISTS pages (
  id BIGSERIAL PRIMARY KEY,
  url TEXT UNIQUE NOT NULL,
  domain TEXT NOT NULL,
  title TEXT,
  status_code INT,
  crawled_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pages_domain ON pages(domain);

CREATE TABLE IF NOT EXISTS links (
  id BIGSERIAL PRIMARY KEY,
  source_url TEXT NOT NULL,
  target_url TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE(source_url, target_url)
);
