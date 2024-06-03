CREATE TABLE IF NOT EXISTS jobs (
    id UUID NOT NULL,
    owner_id UUID NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    responses INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)