-- create
CREATE TABLE IF NOT EXISTS users(
    user_id VARCHAR PRIMARY KEY,
    email VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR NOT NULL,
    refresh_token VARCHAR NOT NULL,
    ip VARCHAR NOT NULL,
    user_id VARCHAR REFERENCES users ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL
);
-- insert
INSERT INTO users VALUES ('92503be1-3a71-4135-a7ea-b42361957c56', 'test@mail.ru');
INSERT INTO users VALUES ('7bc22cf5-b773-48c2-8a1d-b11cf87f22e2', 'test1@mail.ru');
INSERT INTO users VALUES ('6a705be9-dd64-4839-a1ec-76f5094e3935', 'test2@mail.ru');