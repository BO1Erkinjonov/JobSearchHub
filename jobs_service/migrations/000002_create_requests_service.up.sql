CREATE TABLE IF NOT EXISTS requests (
    job_id UUID NOT NULL,
    client_id UUID NOT NULL,
    summary_id INT NOT NULL,
    status_resp VARCHAR(100) CHECK (status_resp IN ('refusal', 'in expectation', 'accepted')) DEFAULT 'in expectation' NOT NULL,
    description_resp TEXT
)