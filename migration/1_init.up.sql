CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    pass_hash TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    birthdate DATE,
    coins INT DEFAULT 1000
);

CREATE TABLE IF NOT EXISTS items
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    price NUMERIC(8,2) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS inventory 
(
    id SERIAL PRIMARY KEY,
    item_id INT,
    username VARCHAR(255),
    quantity INT,
    FOREIGN KEY (username) REFERENCES users(username)
    FOREIGN KEY (item_id) REFERENCES items (item_id)
);

CREATE TABLE IF NOT EXISTS coins_transactions
(
    coins_transaction_id SERIAL PRIMARY KEY,
    from_user VARCHAR(100) NOT NULL,
    to_user VARCHAR(100) NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_user) REFERENCES users (username)
    FOREIGN KEY (to_user) REFERENCES users (username)
);