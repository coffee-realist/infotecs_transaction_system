CREATE TABLE wallets (
    address TEXT PRIMARY KEY,
    balance TEXT NOT NULL
);

CREATE TABLE transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_address TEXT NOT NULL,
    to_address TEXT NOT NULL,
    amount TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_address) REFERENCES wallets(address),
    FOREIGN KEY (to_address) REFERENCES wallets(address)
);
