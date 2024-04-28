GRANT ALL PRIVILEGES ON DATABASE transactions TO docker;

CREATE TABLE accounts (
    id INTEGER PRIMARY KEY,
    balance FLOAT8
);
