CREATE TABLE IF NOT EXISTS Tasks (
    id BIGSERIAL PRIMARY KEY,
    uid BIGSERIAL,
    detail TEXT,
    assignedTo TEXT,
    completeBy timestamptz,
    createdAt timestamptz DEFAULT now(),
    updatedAt timestamptz DEFAULT now()
);