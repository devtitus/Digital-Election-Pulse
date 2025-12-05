-- Enable UUID extension if needed, though we are using integer IDs for simplicity as per plan
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: parties
CREATE TABLE parties (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    leader VARCHAR(255),
    color_hex VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_parties_deleted_at ON parties(deleted_at);

-- Table: sentiment_snapshots
CREATE TABLE sentiment_snapshots (
    id SERIAL PRIMARY KEY,
    party_id INTEGER REFERENCES parties(id),
    score DOUBLE PRECISION,
    key_issue TEXT,
    source_breakdown JSONB, -- Stores JSON like {"yt": 0.5, "news": 0.5}
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Optional: Seed Data to get started
INSERT INTO parties (name, leader, color_hex) VALUES 
('DMK', 'M.K. Stalin', '#dd2e44'),
('AIADMK', 'Edappadi Palaniswami', '#27ae60'),
('TVK', 'Vijay', '#f1c40f'),
('BJP', 'K. Annamalai', '#f39c12'),
('NTK', 'Seeman', '#e74c3c');
