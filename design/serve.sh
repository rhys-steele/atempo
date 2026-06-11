#!/usr/bin/env bash
# Serve the Atempo design system locally — no build step, no node tooling.
# The component cards and the Cockpit kit use Babel's text/babel scripts,
# which need HTTP (not file://). Open http://localhost:4400 once it's up.
set -euo pipefail
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PORT="${1:-4400}"
echo "Atempo design system → http://localhost:${PORT}"
exec python3 -m http.server "${PORT}" -d "${DIR}"
