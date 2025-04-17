#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Set default database URL if not provided
DATABASE_URL=${DATABASE_URL:-"postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"}

docker run --rm \
    -it \
    --network=host \
    -v "$(pwd)/db:/db" \
    -e DATABASE_URL="${DATABASE_URL}" \
    ghcr.io/amacneil/dbmate \
    "$@" 