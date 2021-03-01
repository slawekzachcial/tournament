#!/bin/bash

if ! docker network inspect postgres-net &>/dev/null; then
    docker network create --driver bridge postgres-net
fi

POSTGRES_DATA=$(readlink --canonicalize $(dirname ${BASH_SOURCE[0]})/../postgres)
mkdir -p "${POSTGRES_DATA}"

docker run --rm \
    --name postgres \
    --network postgres-net \
    --publish 5432:5432 \
    --env POSTGRES_PASSWORD=secret \
    --volume "${POSTGRES_DATA}":/var/lib/postgresql/data \
    postgres:alpine
