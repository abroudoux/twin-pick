#!/usr/bin/env bash

set -euo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cat "$DIR/request.json" | "$DIR/../bin/mcp" | jq
