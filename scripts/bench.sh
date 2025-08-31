#!/usr/bin/env bash
set -euo pipefail
BIN=bin/wasmrt
WASM=wasm/hello.wasm

echo "Building runtime and sample..."
go build -o "$BIN" ./cmd/wasmrt
tinygo build -o "$WASM" -target=wasi ./wasm

echo "Cold-start timings (seconds):"
for i in $(seq 1 10); do
  /usr/bin/time -f "%e" "$BIN" create demo$i "$WASM" >/dev/null 2>&1
  /usr/bin/time -f "%e" "$BIN" start demo$i >/dev/null
  "$BIN" delete demo$i >/dev/null 2>&1
done
