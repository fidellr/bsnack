# BSNACK - Point of Sales & Inventory System

BSNACK is a specialized POS and Inventory management system designed for snack businesses. It features a loyalty point system, automated stock management, and a comprehensive sales reporting tool with high-performance caching.

## 🚀 Features

* **Inventory Management:** Track products by type, flavor, size, and manufacturing date.
* **Customer Loyalty:** * Earn 1 point for every **1,000 IDR** spent.
* Redeem points for free products (Small: 200 pts, Medium: 300 pts, Large: 500 pts).


* **Owner Sales Report:** * Aggregated income and products sold.
* Best-selling product detection.
* **New Customer Tracking:** Automatically identifies customers registered within the reporting month.


* **Performance:** Uses **Redis Cache-Aside** strategy for heavy reporting queries.
* **Logging:** Structured JSON logging for production monitoring.

---

## 🛠 Tech Stack

* **Language:** Go 1.21+ (Pure `net/http`)
* **Database:** PostgreSQL (Primary storage)
* **Cache:** Redis (Reporting cache)
* **Log:** `slog` (Structured Logging)
* **Migration:** SQL-based migrations

---

## 📋 Prerequisites

* Go 1.21 or higher
* PostgreSQL 14+
* Redis 6+

---

## ⚙️ Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/fidellr/bsnack.git
cd bsnack

```

### 2. Environment Configuration

Create a `.env` file in the root directory:

```env
APP_ENV=development
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=bsnack_db
REDIS_HOST=localhost:6379
REDIS_PASSWORD=
```

### 3. Database Migration

Run the following commands to set up the schema:

```bash
psql -U postgres -d bsnack_db -f migrations/000001_init_schema.up.sql

```

### 4. Run the Application

```bash
go mod tidy
go run cmd/api/main.go

```

---

## 🔌 API Endpoints

### Products

* `POST /products` - Add new snack inventory.
* `GET /products?date=YYYY-MM-DD` - Get products by manufacturing date.

### Transactions

* `POST /transactions` - Purchase snacks (Supports optional `transaction_date`).
* `GET /transactions?start=YYYY-MM-DD&end=YYYY-MM-DD` - Get Owner Sales Report.

### Redemptions

* `POST /redemptions` - Exchange loyalty points for snacks.


### Customers

* `POST /customers` - Get all the registered customers.

---

## 🧪 Testing

The project includes unit tests for critical business logic (Point calculation, Stock validation, and Cache logic).

```bash
go test ./internal/service/... -v

```

---

## 📖 Business Logic Reference

### Customer Status

A customer is marked as **"New"** in the sales report if their registration date (`created_at`) falls within the same month and year as the transaction being reported.

### Redemption Rules

| Size | Point Cost |
| --- | --- |
| Small | 200 |
| Medium | 300 |
| Large | 500 |
