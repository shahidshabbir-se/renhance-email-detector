CREATE TABLE IF NOT EXISTS results (
    job_id TEXT PRIMARY KEY,
    company TEXT NOT NULL,
    email TEXT,
    fetched_at TIMESTAMP DEFAULT NOW()
);
