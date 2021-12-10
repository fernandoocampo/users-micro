-- Table: public.jobseeker

-- DROP TABLE public.jobseeker;

CREATE TABLE public.jobseeker
(
    id text PRIMARY KEY,
    firstname text COLLATE pg_catalog."default",
    lastname text COLLATE pg_catalog."default",
    city text COLLATE pg_catalog."default",
    skills jsonb
)

TABLESPACE pg_default;

ALTER TABLE public.jobseeker
    OWNER to postgres;
