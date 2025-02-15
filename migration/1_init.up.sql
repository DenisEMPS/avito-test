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
    price INT NOT NULL
);

CREATE TABLE IF NOT EXISTS inventory 
(
    id SERIAL PRIMARY KEY,
    item VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL,
    quantity INT DEFAULT 1,
    FOREIGN KEY (username) REFERENCES users (username),
    FOREIGN KEY (item) REFERENCES items (name)
);

CREATE TABLE IF NOT EXISTS coins_transactions
(
    id SERIAL PRIMARY KEY,
    from_user VARCHAR(100) NOT NULL,
    to_user VARCHAR(100) NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_user) REFERENCES users (username),
    FOREIGN KEY (to_user) REFERENCES users (username)
);

INSERT INTO items (name, price)
VALUES
('t-shirt', 80),
('cup', 20),
('book', 50),
('pen', 10),
('powerbank', 200),
('hoody', 300),
('umbrella', 200),
('socks', 10),
('wallet', 50),
('pink-hoody', 0);
