-- +goose Up
CREATE TABLE users (
    id bigint AUTO_INCREMENT PRIMARY KEY NOT NULL,
    email varchar(128) NOT NULL UNIQUE,
    hashed_password varchar(128) NOT NULL,
    api_key varchar(128) NOT NULL UNIQUE,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

-- +goose Down
DROP TABLE users;