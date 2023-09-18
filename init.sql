-- Create the extension if it doesn't exist (assuming you're using psql)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the jobs table with the specified schema
CREATE TABLE IF NOT EXISTS public.jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payload JSONB,
    reserved_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
    );


ALTER TABLE public.jobs OWNER TO postgres;