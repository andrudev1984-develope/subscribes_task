-- Table: public.subscribe

-- DROP TABLE IF EXISTS public.subscribe;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp"
    SCHEMA public
    VERSION "1.1";


CREATE TABLE IF NOT EXISTS public.subscribe
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    service_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    start_date date NOT NULL,
    end_date date NOT NULL,
    price bigint DEFAULT 0,
    CONSTRAINT subscribe_pkey PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.subscribe
    OWNER to subscriber;