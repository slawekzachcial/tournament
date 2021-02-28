#!/bin/bash

docker run -it --rm \
    --network postgres-net \
    --env PGPASSWORD=secret \
    --volume $(readlink --canonicalize $(dirname ${BASH_SOURCE[0]})/..):/work \
    --workdir /work \
    postgres:alpine \
    psql -h postgres -U postgres "$@"

