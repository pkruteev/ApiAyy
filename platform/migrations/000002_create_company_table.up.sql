-- Set timezone
SET TIMEZONE='Europe/Moscow';

-- Create sequence for company_id if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'company_company_id_seq') THEN
        CREATE SEQUENCE company_company_id_seq;
    END IF;
END $$;

-- Create company_main table
CREATE TABLE IF NOT EXISTS public.company
(
    company_id       integer NOT NULL DEFAULT nextval('company_company_id_seq'::regclass),
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    status_company   character varying(55) NOT NULL,
    company_name     character varying(255),
    inn              character varying(55) NOT NULL,
    kpp              character varying(55),
    ogrn             character varying(55),
    ur_adress        character varying(255),
    mail_adress      character varying(55),
    bank_name        character varying(255),
    bank_bic         character varying(55),
    kor_schet        character varying(55),
    r_schet          character varying(55),
    phone            character varying(25),
    email            character varying(55),
    director         character varying(255),
    CONSTRAINT company_pkey PRIMARY KEY (company_id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.company
    OWNER to admin;
