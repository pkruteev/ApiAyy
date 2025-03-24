-- Set timezone
SET TIMEZONE='Europe/Moscow';

-- Create sequence for company_id if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'companies_company_id_seq') THEN
        CREATE SEQUENCE companies_company_id_seq;
    END IF;
END $$;

-- Create companies table
CREATE TABLE IF NOT EXISTS public.companies
(
    company_id       integer NOT NULL DEFAULT nextval('companies_company_id_seq'::regclass),
    create_company   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    contragent       boolean NOT NULL,
    status           character varying(55) NOT NULL,
    name             character varying(255) NOT NULL,
    inn              character varying(55) NOT NULL,
    kpp              character varying(55),
    ogrn             character varying(55),
    data_ogrn        DATE,
    ogrnip           character varying(55),
    data_ogrnip      DATE,
    ur_address       character varying(255),
    mail_address     character varying(55),
    phone            character varying(25),
    email            character varying(55),
    director         character varying(255),
    CONSTRAINT companies_pkey PRIMARY KEY (company_id)
)
TABLESPACE pg_default;

ALTER TABLE public.companies
OWNER to admin;
