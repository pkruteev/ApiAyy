-- Set timezone
SET TIMEZONE='Europe/Moscow';

-- Create sequence for user_id if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'rights_right_id_seq') THEN
        CREATE SEQUENCE rights_right_id_seq;
    END IF;
END $$;

-- Create users table
CREATE TABLE IF NOT EXISTS public.rights
        (
            rights_id              integer NOT NULL DEFAULT nextval('rights_right_id_seq'::regclass),
            created_at             TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            user_id                integer NOT NULL,
            user_company           integer NOT NULL,
            user_right             character varying(55) NOT NULL,
            CONSTRAINT fk_user_id  FOREIGN KEY (user_id) REFERENCES public.users(user_id),
            CONSTRAINT pk_rights_id   PRIMARY KEY (rights_id)
        )

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.rights
    OWNER to admin;

