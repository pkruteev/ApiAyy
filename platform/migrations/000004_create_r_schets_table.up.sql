SET TIMEZONE='Europe/Moscow';

-- Create sequence for r_schet_id if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'r_schets_r_schet_id_seq') THEN
        CREATE SEQUENCE r_schets_r_schet_id_seq;
    END IF;
END $$;

-- Create r_schets table
CREATE TABLE IF NOT EXISTS public.r_schets
(
    r_schet_id integer NOT NULL DEFAULT nextval('r_schets_r_schet_id_seq'::regclass),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    r_schet character varying(55),
    bank_name character varying(255),
    bank_bic character varying(55),
    kor_schet character varying(55),
    company_id integer, 
    CONSTRAINT r_schets_id PRIMARY KEY (r_schet_id),
    CONSTRAINT fk_companies_id FOREIGN KEY (company_id) REFERENCES public.companies(company_id) 
)
TABLESPACE pg_default;

ALTER TABLE public.r_schets
OWNER to admin;
