\c accountdb;

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    user_id varchar(100) PRIMARY KEY,
    full_name varchar(100) NOT NULL,
    email varchar(100) NOT NULL,
    password varchar(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
