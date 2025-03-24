-- Set timezone
SET TIMEZONE='Europe/Moscow';

-- Create sequence for object_id if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'objects_object_id_seq') THEN
        CREATE SEQUENCE public.objects_object_id_seq;
    END IF;
END $$;

-- Create objects table
CREATE TABLE IF NOT EXISTS public.objects
(
    object_id       integer NOT NULL DEFAULT nextval('public.objects_object_id_seq'::regclass),
    created_ob      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    typereal        character varying(55),
    city            character varying(255),
    street          character varying(255),
    house           character varying(55),
    flat            character varying(55),
    square          character varying(55),
    floor           character varying(55),
    company_id      integer, 
    CONSTRAINT pk_object_id PRIMARY KEY (object_id),
    CONSTRAINT fk_company_id FOREIGN KEY (company_id) REFERENCES public.companies(company_id) 
)
TABLESPACE pg_default;

ALTER TABLE public.objects
OWNER to admin;
