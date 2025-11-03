#!/usr/bin/env bash

set -euo pipefail

BASE_URL="http://localhost:8080/api"

get() {
    local endpoint="$1"

    echo "➡️ GET $endpoint"

    response=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL$endpoint")

    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)

    echo "Status: $http_code"
    echo "Body: $body"

    if [[ "$http_code" != "200" && "$http_code" != "201" ]]; then
        echo "Error: $body"
    fi

    echo ""
}

endpoints=(
    "/v1/pick?usernames=abroudoux,potatoze&limit=1"
    "/v1/pick?usernames=abroudoux,potatoze&genres=action&limit=10"
)

for ep in "${endpoints[@]}"; do
    get "$ep"
done
