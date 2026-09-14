#!/usr/bin/env bash
# Thin wrapper around gen_stats.py
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"
exec python3 src/examples/gen_stats.py "$@"
