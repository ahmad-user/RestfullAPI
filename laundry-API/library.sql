-- Tabel Customer
CREATE TABLE tbl_customer(
    id SERIAL PRIMARY KEY,
    nama_customer VARCHAR(100),
    no_hp VARCHAR(20),
    address VARCHAR(200)
);

INSERT INTO tbl_customer (nama_customer, no_hp, address) VALUES 
    ('Jessica', '021', 'Jakarta'),
    ('Mirna', '022', 'Jakarta');

-- Tabel Products
CREATE TABLE tbl_products(
    id SERIAL PRIMARY KEY,
    Name VARCHAR(100) NOT NULL,
    Price INT NOT NULL,
    Unit VARCHAR(20) NOT NULL
);

INSERT INTO tbl_products (Name, Price, Unit) VALUES 
    ('joko', 6000, “KG),
    ('ahmad', 5000, “Buah”);

-- Tabel Employees
CREATE TABLE tbl_employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    address TEXT
);

INSERT INTO tbl_employees (name, phone_number, address) VALUES 
    ('joko', '23', 'sangiang'),
    ('ahmad', '24', 'jakarta');

-- Tabel Transactions
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    bill_date DATE NOT NULL,
    entry_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finish_date TIMESTAMP,
    employee_id INT NOT NULL,
    customer_id INT NOT NULL,
    total_bill INT NOT NULL,
    FOREIGN KEY (employee_id) REFERENCES tbl_employees(id),
    FOREIGN KEY (customer_id) REFERENCES tbl_customer(id)
);

INSERT INTO transactions (bill_date, entry_date, finish_date, employee_id, nama_employee, customer_id, total_bill) VALUES 
    ('2022-04-27', '2022-04-27 12:00:00', '2022-04-27 15:00:00', 1, 'joko', 1, 200000);

-- Tabel Bill Details
CREATE TABLE bill_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT NOT NULL,
    product_id INT NOT NULL,
    qty INT NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (product_id) REFERENCES tbl_products(id)
);

INSERT INTO bill_details (transaction_id, product_id, qty) VALUES 
    (1, 1, 3),
    (1, 2, 1);
