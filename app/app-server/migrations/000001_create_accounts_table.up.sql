CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    balance DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_accounts_email ON accounts(email);
CREATE INDEX idx_accounts_created_at ON accounts(created_at);

DROP TABLE IF EXISTS accounts;
