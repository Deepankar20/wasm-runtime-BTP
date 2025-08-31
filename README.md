# wasm-runtime (WASM-native "container" runtime in Go)

This is a minimal WASM-specific runtime written in Go. It provides a container-like CLI (`create/start/delete`) but executes WebAssembly modules directly via a pure Go engine (wazero) instead of spawning Linux containers.

## Quickstart
```bash
go mod tidy
make build
make wasm             # requires tinygo
bin/wasmrt create demo wasm/hello.wasm
bin/wasmrt start demo
bin/wasmrt delete demo
