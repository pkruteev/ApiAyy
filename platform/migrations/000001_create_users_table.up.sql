-- Set timezone
SET TIMEZONE='Europe/Moscow';

-- Create sequence for user_id if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'users_user_id_seq') THEN
        CREATE SEQUENCE users_user_id_seq;
    END IF;
END $$;

-- Create users table
CREATE TABLE IF NOT EXISTS public.users
(
    user_id          integer NOT NULL DEFAULT nextval('users_user_id_seq'::regclass),
    created_user     TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    bd_used          integer,
    first_name       character varying(50) NOT NULL,
    patronymic_name  character varying(50),
    last_name        character varying(50),
    user_email       character varying(50) NOT NULL UNIQUE,
    user_phone       character varying(20),
    password         character varying(255) NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (user_id)
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to admin;
