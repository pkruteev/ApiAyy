-- Set timezone
SET TIMEZONE='Europe/Moscow';

-- Create user_role_enum type if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role_enum') THEN
        CREATE TYPE user_role_enum AS ENUM ('member', 'admin', 'director', 'manager');
    END IF;
END $$;

-- Create sequence for user_id if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'rights_right_id_seq') THEN
        CREATE SEQUENCE rights_right_id_seq;
    END IF;
END $$;

-- Create users table
CREATE TABLE IF NOT EXISTS public.rights
        (
            rights_id              integer NOT NULL DEFAULT nextval('rights_right_id_seq'::regclass),
            user_id                integer NOT NULL,
            user_bd                integer NOT NULL,
            holding                character varying(55),
            user_rights            user_role_enum NOT NULL,
            CONSTRAINT fk_user_id  FOREIGN KEY (user_id) REFERENCES public.users(user_id),
            CONSTRAINT pk_rights_id   PRIMARY KEY (rights_id)
        )

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.rights
    OWNER to admin;

    -- Создаем индекс для ускорения поиска по user_id
-- CREATE INDEX idx_rights_user_id ON public.rights (user_id);
-- Создание индекса для ускорения поиска по user_id, если таблица существует
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'rights') THEN
        CREATE INDEX idx_rights_user_id ON public.rights (user_id);
    END IF;
END $$;
