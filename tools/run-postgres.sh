#!/bin/bash

if ! docker network inspect postgres-net &>/dev/null; then
    docker network create --driver bridge postgres-net
fi

docker run --rm \
    --name postgres \
    --network postgres-net \
    --publish 5432:5432 \
    --env POSTGRES_PASSWORD=secret \
    --volume $(readlink --canonicalize $(dirname ${BASH_SOURCE[0]})/../postgres):/var/lib/postgresql/data \
    postgres:alpine
