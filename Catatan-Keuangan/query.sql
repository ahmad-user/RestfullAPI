select * from expenses;
select * from users;
alter table expenses drop column user_id;
alter table expenses add column updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE transaction_type AS ENUM ('CREDIT', 'DEBIT');
CREATE TABLE expenses (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    date DATE NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    transaction_type transaction_type,
    balance DOUBLE PRECISION NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);
INSERT INTO expenses (date, amount, transaction_type, balance, description)
VALUES 
('2024-05-01', 150.75, 'CREDIT', 1500.00, 'Salary payment'),
('2024-05-01', 150.75, 'CREDIT', 1500.00, 'Salary payment');

CREATE TABLE users (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  username VARCHAR(100) unique,
  password VARCHAR(100),
  role VARCHAR(30),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
);
INSERT INTO users (username, password, role, created_at, updated_at)
VALUES ('joko2', 'pass', 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('joko3', 'pass', 'user', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

ALTER TABLE expenses ADD COLUMN user_id uuid;
ALTER TABLE expenses ADD FOREIGN KEY (user_id) REFERENCES users(id)