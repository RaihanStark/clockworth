#!/bin/bash
set -e

echo "Running database migrations..."

# Check if dbmate is installed
if ! command -v dbmate &> /dev/null; then
    echo "dbmate not found, installing..."
    go install github.com/amacneil/dbmate/v2@latest
fi

# Run migrations
dbmate up

echo "Migrations completed successfully!"
