GRANT ALL PRIVILEGES ON DATABASE transactions TO docker;

CREATE TABLE accounts (
    id INTEGER PRIMARY KEY,
    balance FLOAT8
);


CREATE TABLE transaction_ledger (
    id INTEGER PRIMARY KEY,
    source_account_id INTEGER,
    destination_account_id INTEGER,
    amount FLOAT8
);
