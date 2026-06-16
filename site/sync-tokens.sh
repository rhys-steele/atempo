#!/usr/bin/env bash
# Copy the design system's tokens into site/css/tokens/ so site/ stays a
# self-contained, deployable folder. design/tokens/ is the source of truth —
# never hand-edit the copies. Re-run whenever tokens change; commit the output.
set -euo pipefail
here="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cp "$here"/../design/tokens/*.css "$here"/css/tokens/
echo "synced tokens → site/css/tokens/"
