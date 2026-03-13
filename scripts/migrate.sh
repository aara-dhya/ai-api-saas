#!/bin/bash

set -e

if [ -z "$DATABASE_URL" ]; then
  echo "DATABASE_URL is not set"
  exit 1
fi

echo "Running database migrations..."

for file in backend/migrations/*.sql
do
  echo "Applying $file"
  psql "$DATABASE_URL" -f "$file"
done

echo "Migrations completed"