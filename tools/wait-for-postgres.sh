#!/bin/sh

set -e

# until PGPASSWORD="$POSTGRES_PASSWORD" psql -h "$host" -U "$POSTGRES_USER" -l | grep -q "$POSTGRES_DB"; do
until psql "$DB_URL" -l | grep -q tournament; do
    >&2 echo 'Database unavailable - sleeping'
    sleep 5
done

>&2 echo 'Database available - starting application'
exec "$@"
