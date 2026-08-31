#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

mkdir -p "$ROOT_DIR/dist"

GOOS=wasip1 GOARCH=wasm go build \
  -o "$ROOT_DIR/dist/plugin.wasm" \
  "$ROOT_DIR/cmd/plugin/main.go"

echo "Built dist/plugin.wasm"
