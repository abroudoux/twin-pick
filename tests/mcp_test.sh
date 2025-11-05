#!/usr/bin/env bash

set -euo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

cat "$DIR/request_pick.json" | "$DIR/../bin/twinpick-mcp" | jq
cat "$DIR/request_spot.json" | "$DIR/../bin/twinpick-mcp" | jq
