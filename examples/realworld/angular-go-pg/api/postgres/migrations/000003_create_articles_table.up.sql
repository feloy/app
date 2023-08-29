BEGIN;

CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    description TEXT,
    slug VARCHAR(255) NOT NULL UNIQUE,
    author_id  INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_author FOREIGN KEY(author_id) REFERENCES users(id)
);

COMMIT;