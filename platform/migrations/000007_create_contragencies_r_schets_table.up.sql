-- Set timezone
SET TIMEZONE='Europe/Moscow';

-- Create sequence for contr_schet_id if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'contr_schet_contr_schet_id_seq') THEN
        CREATE SEQUENCE contr_schet_contr_schet_id_seq;
    END IF;
END $$;

-- Create contr_schet table
CREATE TABLE IF NOT EXISTS public.contr_schet
(
    contr_schet_id       integer NOT NULL DEFAULT nextval('contr_schet_contr_schet_id_seq'::regclass),
    created_at           TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    contr_schet          character varying(55),
    bank_name            character varying(255),
    bank_bic             character varying(55),
    kor_schet            character varying(55),
    contragency_id       integer, 
    CONSTRAINT pkey_contr_schet_id PRIMARY KEY (contr_schet_id),
    CONSTRAINT fk_contragency_id FOREIGN KEY (contragency_id) REFERENCES public.contragencies(contragency_id) 
)
TABLESPACE pg_default;

ALTER TABLE public.contr_schet
    OWNER TO admin;
