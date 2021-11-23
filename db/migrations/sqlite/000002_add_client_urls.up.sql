CREATE TABLE redirect_url (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id INTEGER NOT NULL,
    url VARCHAR NOT NULL,
    FOREIGN KEY(client_id) REFERENCES client(id)
);