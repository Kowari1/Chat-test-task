-- +goose Up
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL
    CHECK (LENGTH(TRIM(title)) > 0 AND LENGTH(TRIM(title)) <= 200),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE chats;
