GRANT ALL PRIVILEGES ON DATABASE transactions TO docker;

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    balance FLOAT8,
    version INTEGER
);


CREATE TABLE transactions(
    id SERIAL PRIMARY KEY,
    source_account_id INTEGER,
    destination_account_id INTEGER,
    amount FLOAT8
);
