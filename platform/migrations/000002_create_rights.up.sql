-- Устанавливаем временную зону
SET TIMEZONE='Europe/Moscow';

-- Создаем тип user_role_enum, если он не существует
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role_enum') THEN
        CREATE TYPE user_role_enum AS ENUM ('member', 'admin', 'director', 'manager');
    END IF;
END $$;

-- Создаем последовательность для rights_id, если она не существует
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'rights_right_id_seq') THEN
        CREATE SEQUENCE rights_right_id_seq;
    END IF;
END $$;

-- Создаем таблицу rights, если она не существует
CREATE TABLE IF NOT EXISTS public.rights
(
    rights_id         integer NOT NULL DEFAULT nextval('rights_right_id_seq'::regclass),
    user_id           integer NOT NULL,
    created_rights    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    user_bd           character varying(55) NOT NULL,
    holding           character varying(55),
    user_role         user_role_enum NOT NULL,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES public.users(user_id),
    CONSTRAINT pk_rights_id PRIMARY KEY (rights_id)
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.rights
    OWNER TO admin;

-- Создаем обычный индекс на user_id для ускорения поиска, если он не существует
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_indexes 
        WHERE tablename = 'rights' 
        AND indexname = 'idx_rights_user_id'
    ) THEN
        CREATE INDEX idx_rights_user_id ON public.rights (user_id);
    END IF;
END $$;