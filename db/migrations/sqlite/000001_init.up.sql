CREATE TABLE user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email VARCHAR UNIQUE NOT NULL,
    hash VARCHAR NOT NULL
);

CREATE TABLE client (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id VARCHAR UNIQUE NOT NULL,
    client_secret VARCHAR NOT NULL
);