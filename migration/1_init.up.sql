CREATE TABLE IF NOT EXISTS users
(
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    pass_hash TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    birthdate DATE,
    coins INT DEFAULT 1000
);

CREATE TABLE IF NOT EXISTS products
(
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    price NUMERIC(8,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS users_products
(
    users_products_id SERIAL PRIMARY KEY,
    user_id INT,
    product_id INT,
    FOREIGN KEY (user_id) REFERENCES users (user_id),
    FOREIGN KEY (product_id) REFERENCES products (product_id)
);

CREATE TABLE IF NOT EXISTS coins_transactions
(
    coins_transaction_id SERIAL PRIMARY KEY,
    from_user_id INT,
    to_user_id INT,
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending',
    FOREIGN KEY (from_user_id) REFERENCES users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (to_user_id) REFERENCES users (user_id) ON DELETE CASCADE
);
