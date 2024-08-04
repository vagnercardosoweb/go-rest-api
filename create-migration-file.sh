#!/usr/bin/env sh

set -e

if [ -z "$1" ]; then
  echo "Please provide a name for the migration"
  echo "Example: ./create-migration-file.sh migration_name"
  exit 1
fi

prefix=$(date +%Y%m%d%H%M%S)
filename=$(echo "$1" | sed 's/[^a-zA-Z0-9]/_/g' | tr '[:upper:]' '[:lower:]')

touch "migrations/${prefix}_$filename.up.sql"
echo "✅ Created migrations/${prefix}_$filename.up.sql"

touch "migrations/${prefix}_$filename.down.sql"
echo "✅ Created migrations/${prefix}_$filename.down.sql"
