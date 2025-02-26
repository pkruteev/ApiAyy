SET TIMEZONE='Europe/Moscow';

-- Create sequence for contragency_id if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'contragencies_contragency_id_seq') THEN
        CREATE SEQUENCE contragencies_contragency_id_seq;
    END IF;
END $$;

-- Create contragencies table
CREATE TABLE IF NOT EXISTS public.contragencies
(
    contragency_id integer NOT NULL DEFAULT nextval('contragencies_contragency_id_seq'::regclass),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    contragency_status character varying(55) NOT NULL,
    contragency_name character varying(255),
    inn character varying(55) NOT NULL,
    kpp character varying(55),
    ogrn character varying(55),
    ur_adress character varying(255),
    mail_adress character varying(55),
    phone character varying(25),
    email character varying(55),
    director character varying(255),
    CONSTRAINT contragency_id_pkey PRIMARY KEY (contragency_id)
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.contragencies
OWNER to admin;
