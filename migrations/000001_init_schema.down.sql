-- drop tables in reverse order of creation to avoid Foreign Key violations
DROP TABLE IF EXISTS redemptions;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS customers;
DROP TABLE IF EXISTS products;

DROP EXTENSION IF EXISTS "uuid-ossp";