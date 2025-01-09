-- Set timezone
SET TIMEZONE='Europe/Moscow';

-- Create sequence for contract_id if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'contracts_contract_id_seq') THEN
        CREATE SEQUENCE contracts_contract_id_seq;
    END IF;
END $$;

-- Create contracts table
CREATE TABLE IF NOT EXISTS public.contracts

        (
            contract_id          integer NOT NULL DEFAULT nextval('contracts_contract_id_seq'::regclass),
            created_at           TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            data_start_contract          character varying(55),
            data_end_contract            character varying(255),
            data_start_rent             character varying(55),
            data_end_rent            character varying(55),
            object_id               integer,
            company_id            integer,
            contragency_id         integer, 
            service_contract            bool,
            rent_pay                     character varying(255),
            rent_pre_pay                 character varying(255),

            CONSTRAINT pkey_contract_id PRIMARY KEY (contract_id),
            CONSTRAINT fk_object_id FOREIGN KEY (object_id) REFERENCES public.objects(object_id),
            CONSTRAINT fk_company_id FOREIGN KEY (company_id) REFERENCES public.companies(company_id),
            CONSTRAINT fk_contragency_id FOREIGN KEY (contragency_id) REFERENCES public.contragencies(contragency_id),
        )
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.contracts
    OWNER to admin;
