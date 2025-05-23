#!/bin/sh

PORT=$(echo "$SERVER_METRICS_ADDRESS" | cut -d ':' -f2)

PORT=${PORT:-5082}

URL="http://127.0.0.1:$PORT/health"

wget --no-verbose --spider "$URL" || exit 1
