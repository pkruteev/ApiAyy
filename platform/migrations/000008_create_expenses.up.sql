-- Set timezone
SET TIMEZONE='Europe/Moscow';

-- Create sequence for contract_id if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'statements_statement_id_seq') THEN
        CREATE SEQUENCE statements_statement_id_seq;
    END IF;
END $$;

-- Create contracts table
CREATE TABLE IF NOT EXISTS public.statements
(
    statement_id         integer NOT NULL DEFAULT nextval('statements_statement_id_seq'::regclass),
    created_statement    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    rschet_id            integer,
    r_schet              character varying(55),
    date_transaction     DATE,
    company_id           integer,
    contragent_id        integer,
    balance_begin_day    character varying(55),
    kredit               character varying(55),
    debit                character varying(55),
    balance_end_day      character varying(55),
    basis_payment        character varying(55),
    author_id            integer,

    CONSTRAINT pk_statement_id PRIMARY KEY (statement_id),
    CONSTRAINT fk_rschet_id FOREIGN KEY (rschet_id) REFERENCES public.accounts(rschet_id) 
    CONSTRAINT fk_r_schet FOREIGN KEY (r_schet) REFERENCES public.accounts(r_schet) 
    CONSTRAINT fk_company_id FOREIGN KEY (company_id) REFERENCES public.companies(company_id) 
    CONSTRAINT fk_contragent_id FOREIGN KEY (contragent_id) REFERENCES public.companies(company_id) 
    CONSTRAINT fk_author_id FOREIGN KEY (author_id) REFERENCES public.users(user_id) 
)

TABLESPACE pg_default;

ALTER TABLE public.statements
    OWNER TO admin;
