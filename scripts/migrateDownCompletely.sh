#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd migrations/postgres
goose postgres $DATABASE_URL down-to 000