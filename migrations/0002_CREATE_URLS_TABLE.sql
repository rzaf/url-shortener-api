-- +goose Up
CREATE TABLE urls (
    `id` bigint AUTO_INCREMENT PRIMARY KEY NOT NULL,
    `url` varchar(128) NOT NULL,
    `user_id` bigint NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT NOW(),
    `updated_at` TIMESTAMP NULL,
    CONSTRAINT UNIQUE (`user_id`,`url`),
    CONSTRAINT FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE urls;
