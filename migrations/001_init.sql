-- 001_init.sql

-- Тип статуса PR
CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

-- Таблица команд
CREATE TABLE teams (
    name TEXT PRIMARY KEY
);

-- Таблица пользователей
CREATE TABLE users (
    id        TEXT PRIMARY KEY,
    username  TEXT NOT NULL,
    team_name TEXT NOT NULL REFERENCES teams(name) ON DELETE RESTRICT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_users_team_name ON users(team_name);

-- Таблица pull requests
CREATE TABLE pull_requests (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    author_id  TEXT NOT NULL REFERENCES users(id),
    status     pr_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    merged_at  TIMESTAMPTZ
);

CREATE INDEX idx_pull_requests_author ON pull_requests(author_id);

-- Таблица ревьюверов PR
CREATE TABLE pull_request_reviewers (
    pr_id       TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id TEXT NOT NULL REFERENCES users(id),
    PRIMARY KEY (pr_id, reviewer_id)
);

CREATE INDEX idx_pr_reviewers_reviewer ON pull_request_reviewers(reviewer_id);
