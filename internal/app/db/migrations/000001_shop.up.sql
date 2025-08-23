CREATE TABLE  IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL ,
    password TEXT NOT NULL, 
    balance INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS transfers (
    id SERIAL PRIMARY KEY,
    sender INT NOT NULL,
    receiver INT NOT NULL,
    amount INT NOT NULL CHECK (amount > 0),
    made_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (sender <> receiver)
);

CREATE TABLE  IF NOT EXISTS merch (
    item_id SERIAL PRIMARY KEY,
    item_name VARCHAR(100) UNIQUE NOT NULL,
    cost INT NOT NULL 
);

CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    userID VARCHAR(50) REFERENCES users(id) ON DELETE CASCADE,
    itemID INT REFERENCES merch(item_id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity >= 0),
    UNIQUE (userID, itemID) 
);

INSERT INTO merch (item_name, cost) VALUES
('t-shirt', 80),
('cup', 20),
('book', 50),
('pen', 10),
('powerbank', 200),
('hoody', 300),
('umbrella', 200),
('socks', 10),
('wallet', 50),
('pink-hoody', 500);
