-- Set timezone
SET TIMEZONE='Europe/Moscow';

-- Create sequence for user_id if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'companies_company_id_seq') THEN
        CREATE SEQUENCE companies_company_id_seq;
    END IF;
END $$;

-- Create users table
CREATE TABLE IF NOT EXISTS public.companies
        (
            company_id       integer NOT NULL DEFAULT nextval('companies_company_id_seq'::regclass),
            created_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            status_company   character varying(55) NOT NULL,
            company_name     character varying(255),
            inn              character varying(55) NOT NULL,
            kpp              character varying(55),
            ogrn             character varying(55),
            ur_adress        character varying(255),
            mail_adress      character varying(55),
            phone            character varying(25),
            email            character varying(55),
            director         character varying(255),
            CONSTRAINT companies_pkey PRIMARY KEY (company_id)
        )

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.companies
    OWNER to admin;
