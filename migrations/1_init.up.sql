CREATE TABLE IF NOT EXISTS users
(
    id        BINARY(16) PRIMARY KEY,
    email     VARCHAR(60) NOT NULL UNIQUE,
    pass_hash BLOB NOT NULL
);
CREATE INDEX idx_email ON users(email);

CREATE TABLE IF NOT EXISTS apps
(
    id     INT PRIMARY KEY,
    name   VARCHAR(30) NOT NULL UNIQUE,
    secret VARCHAR(60) NOT NULL UNIQUE
);
CREATE INDEX idx_binary_id ON apps(id);