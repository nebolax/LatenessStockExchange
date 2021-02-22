CREATE TABLE 'users' (
    'id' INTEGER PRIMARY KEY AUTOINCREMENT,
    'username' VARCHAR (255) UNIQUE NOT NULL,
    'email' VARCHAR (255) UNIQUE NOT NULL,
    'password_hash' VARCHAR (255) NOT NULL,
    'password_salt' VARCHAR (255) NOT NULL
);

CREATE TABLE 'stocks' (
    'id' INTEGER PRIMARY KEY AUTOINCREMENT,
    'name' VARCHAR (255) UNIQUE NOT NULL,
    'user_id' INTEGER NOT NULL,
    'come_in_time' DATETIME NULL
);

CREATE TABLE 'user_stock_ownerships' (
    'user_id' INTEGER NOT NULL,
    'stock_id' INTEGER NOT NULL,
    'amount' INTEGER NOT NULL
);

CREATE TABLE 'price_logs' (
    'stock_id' INTEGER NOT NULL,
    'price' FLOAT(10) NOT NULL,
    'timestamp' DATETIME NULL
);

CREATE TABLE 'transaction_logs' (
    'user_id' INTEGER NOT NULL,
    'stock_id' INTEGER NOT NULL,
    'amount' INTEGER NOT NULL,
    'money_spent' INTEGER,
    'type' INTEGER,
    'timestamp' DATETIME NULL
);

CREATE TABLE 'comes_in' (
    'certifier_id' INTEGER NOT NULL,
    'stock_id' INTEGER NOT NULL,
    'timestamp' DATETIME NULL
);

CREATE TABLE 'event_logs' (
    'event_type' INTEGER,
    'user_id' INTEGER NOT NULL,
    'stock_id' INTEGER NOT NULL,
    'timestamp' DATETIME NULL
);
