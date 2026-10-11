#!/bin/sh
set -e
cd "$(dirname "$0")"
npm run build:embed
echo "synced to internal/server/admin-ui"
