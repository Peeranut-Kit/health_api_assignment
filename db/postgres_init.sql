-- Create a new database
-- CREATE DATABASE healthcare;

-- Switch to the new database
-- \c healthcare;
-- This command is specific to the psql command-line interface and is not valid in an SQL script that PostgreSQL executes during container initialization.

-- Create a "hospital" table
CREATE TABLE IF NOT EXISTS hospitals (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

-- Create a "patient" table
CREATE TABLE IF NOT EXISTS patients (
    id SERIAL PRIMARY KEY,
    first_name_th VARCHAR(255),
    middle_name_th VARCHAR(255),
    last_name_th VARCHAR(255),
    first_name_en VARCHAR(255),
    middle_name_en VARCHAR(255),
    last_name_en VARCHAR(255),
    date_of_birth DATE NOT NULL,
    patient_hn VARCHAR(50) NOT NULL UNIQUE,
    national_id VARCHAR(50) NOT NULL UNIQUE,
    passport_id VARCHAR(50) NOT NULL UNIQUE,
    phone_number VARCHAR(50) NOT NULL,
    email VARCHAR(255) UNIQUE,
    gender CHAR(1),
    hospital_id INT REFERENCES hospitals(id) -- Foreign key
);

-- Create a "staff" table
CREATE TABLE IF NOT EXISTS staffs (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    hospital_id INT REFERENCES hospitals(id) -- Foreign key
);