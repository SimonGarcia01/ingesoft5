#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "postgres" <<-EOSQL
    CREATE TABLE IF NOT EXISTS songs (
        songid serial PRIMARY KEY,
        name TEXT,
        singer TEXT,
        genre TEXT
    );
EOSQL
