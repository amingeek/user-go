-- migrations/init.sql
CREATE TABLE IF NOT EXISTS users (
                                     phone VARCHAR(255) PRIMARY KEY,
    registration_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
