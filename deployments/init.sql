CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE subscribes
(
    id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id           TEXT NOT NULL,
    tag TEXT NOT NULL,
    readen_articles text[]
);