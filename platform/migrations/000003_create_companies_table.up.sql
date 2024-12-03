-- Создание таблицы companies, если она не существует
CREATE TABLE IF NOT EXISTS public.companies
(
    company_id serial PRIMARY KEY, -- Используйте тип serial для автоматической генерации последовательностей
    first_name character varying(50) NOT NULL,
    patronymic_name character varying(50)
)
TABLESPACE pg_default;

-- Установка владельца таблицы
ALTER TABLE public.companies
    OWNER TO admin;
