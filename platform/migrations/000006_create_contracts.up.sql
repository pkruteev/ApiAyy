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
    contract_id           integer NOT NULL DEFAULT nextval('contracts_contract_id_seq'::regclass),
    created_contract      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    contract_number       character varying(255),
    date_start            DATE,
    date_end              DATE,
    date_start_rent       DATE,
    date_end_rent         DATE,
    object_id             integer,
    company_id            integer,
    contragent_id         integer, 
    rent_pay              character varying(255),
    day_payment_rent      INTEGER CHECK (day_payment_rent BETWEEN 1 AND 31),
    rent_pre_pay          character varying(255),
    date_rent_prepay      DATE,
    utility_bills         boolean,

    CONSTRAINT pkey_contract_id PRIMARY KEY (contract_id),
    CONSTRAINT fk_object_id FOREIGN KEY (object_id) REFERENCES public.objects(object_id),
    CONSTRAINT fk_company_id FOREIGN KEY (company_id) REFERENCES public.companies(company_id),
    CONSTRAINT fk_contragency_id FOREIGN KEY (contragency_id) REFERENCES public.contragencies(contragency_id)
)

TABLESPACE pg_default;

ALTER TABLE public.contracts
    OWNER TO admin;
