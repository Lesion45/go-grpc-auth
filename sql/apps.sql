CREATE SCHEMA IF NOT EXISTS apps_schema;

CREATE TABLE IF NOT EXISTS app(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT,
    secret TEXT
);