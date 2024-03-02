#!/usr/bin/env sh

set -e

if [ -z "$1" ]
then
  echo "Please provide a name for the migration"
  echo "Example: make create_migration name=..."
  exit 1
fi

prefix=$(date +%Y%m%d%H%M%S)

touch "migrations/${prefix}_$1.up.sql"
echo "Created migrations/${prefix}_$1.up.sql"

touch "migrations/${prefix}_$1.down.sql"
echo "Created migrations/${prefix}_$1.down.sql"
