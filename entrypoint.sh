#!/bin/bash
set -e

echo "Waiting for PostgreSQL..."
timeout 30 sh -c 'until pg_isready -h postgres_db -p 5432 -U admin; do
  echo "PostgreSQL is unavailable - sleeping..."
  sleep 2
done' || {
  echo "PostgreSQL is not ready after 30 seconds!";
  exit 1;
}

echo "Running migrations..."
goose -dir /app/migrations postgres "user=$PG_USER password=$PG_PASS dbname=stocklk host=postgres_db port=5432 sslmode=disable" up || {
  echo "Migration failed!";
  exit 1;
}

echo "Starting application..."
exec /app/main