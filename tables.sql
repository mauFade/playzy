-- users
CREATE TABLE users (id UUID PRIMARY KEY, name VARCHAR NOT NULL, email VARCHAR NOT NULL, phone VARCHAR NOT NULL, password VARCHAR NOT NULL, gamertag VARCHAR NOT NULL, is_deleted BOOLEAN NOT NULL, deleted_at TIMESTAMP NULL, updated_at TIMESTAMP NOT NULL, created_at TIMESTAMP NOT NULL);

-- sessions
CREATE TABLE sessions (id UUID PRIMARY KEY, game VARCHAR NOT NULL, user_id UUID NOT NULL, objective VARCHAR NOT NULL, rank VARCHAR NULL, is_ranked BOOLEAN NOT NULL, updated_at TIMESTAMP NOT NULL, created_at TIMESTAMP NOT NULL, CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE);

