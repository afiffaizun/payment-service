CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    balance BIGINT DEFAULT 0 CHECK (balance >= 0),
    version INT DEFAULT 1
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    reference_id VARCHAR(100) UNIQUE,
    sender_id UUID REFERENCES users(id),
    receiver_id UUID REFERENCES users(id),
    amount BIGINT NOT NULL,
    status VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Data Dummy untuk Testing Awal
INSERT INTO users (id, username) VALUES 
('11111111-1111-1111-1111-111111111111', 'alice'),
('22222222-2222-2222-2222-222222222222', 'bob');

INSERT INTO wallets (id, user_id, balance) VALUES 
('aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 1000000),
('bbbb2222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '22222222-2222-2222-2222-222222222222', 50000);