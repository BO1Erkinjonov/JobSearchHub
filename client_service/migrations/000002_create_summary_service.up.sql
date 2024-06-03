CREATE TABLE IF NOT EXISTS summary (
    id SERIAL ,
    owner_id UUID NOT NULL ,
    skills TEXT NOT NULL ,
    bio TEXT NOT NULL ,
    languages TEXT NOT NULL
)