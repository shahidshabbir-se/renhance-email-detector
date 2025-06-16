-- +migrate Up
CREATE TABLE jobs (
    id TEXT PRIMARY KEY,
    domain TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE companies (
    id SERIAL PRIMARY KEY,
    domain TEXT UNIQUE NOT NULL,
    organization TEXT,
    description TEXT,
    country TEXT,
    city TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE contacts (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL REFERENCES companies (id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    position TEXT,
    department TEXT,
    linkedin_url TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE job_results (
    job_id TEXT NOT NULL REFERENCES jobs (id) ON DELETE CASCADE,
    company_id INTEGER NOT NULL REFERENCES companies (id) ON DELETE CASCADE,
    PRIMARY KEY (job_id, company_id)
);
