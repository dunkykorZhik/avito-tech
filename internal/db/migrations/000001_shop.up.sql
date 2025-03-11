CREATE TABLE  IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL ,
    password TEXT NOT NULL, 
    balance INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS transfers (
    id SERIAL PRIMARY KEY,
    sender VARCHAR(50) NOT NULL,
    receiver VARCHAR(50) NOT NULL,
    amount INT NOT NULL CHECK (amount > 0),
    made_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender) REFERENCES users(username) ON DELETE CASCADE,
    FOREIGN KEY (receiver) REFERENCES users(username) ON DELETE CASCADE,
    CHECK (sender <> receiver)
);

CREATE TABLE  IF NOT EXISTS merch (
    item_id SERIAL PRIMARY KEY,
    item_name VARCHAR(100) UNIQUE NOT NULL,
    cost INT NOT NULL 
);

CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    username INT REFERENCES users(username) ON DELETE CASCADE,
    item_name  VARCHAR(100) REFERENCES merch(item_name) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity >= 0),
    UNIQUE (user_id, item_id) 
);
