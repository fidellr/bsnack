-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Create Products Table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(100) NOT NULL,
    flavor VARCHAR(50) NOT NULL,
    size VARCHAR(20) NOT NULL,
    price NUMERIC(15, 2) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    manufacturing_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create Customers Table
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    points INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Create Transactions Table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id INT NOT NULL REFERENCES customers(id),
    product_id INT NOT NULL REFERENCES products(id),
    quantity INT NOT NULL,
    total_price NUMERIC(15, 2) NOT NULL,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Create Performance Indexes
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
CREATE INDEX idx_products_manuf_date ON products(manufacturing_date);
CREATE INDEX idx_customers_name ON customers(name);