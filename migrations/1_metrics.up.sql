CREATE TABLE IF NOT EXISTS metrics (
    ts TIMESTAMP WITH TIME ZONE PRIMARY KEY,
    cpu_load FLOAT NOT NULL,
    concurrency INT NOT NULL
);
