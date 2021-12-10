#!/bin/bash
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d "$POSTGRES_DB"  <<-EOSQL
     create schema if not exists $SCHEMA;
     CREATE TABLE $SCHEMA.jobseeker
     (
        id text PRIMARY KEY,
        firstname text COLLATE pg_catalog."default",
        lastname text COLLATE pg_catalog."default",
        city text COLLATE pg_catalog."default",
        skills jsonb
     );
     ALTER TABLE $SCHEMA.jobseeker
        OWNER to postgres;
EOSQL