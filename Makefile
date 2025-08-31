BIN=bin/wasmrt
WASM=wasm/hello.wasm

all: build

build:
	go build -o $(BIN) ./cmd/wasmrt

wasm:
	# TinyGo build (install tinygo first: https://tinygo.org/getting-started/)
	tinygo build -o $(WASM) -target=wasi ./wasm

run: build wasm
	$(BIN) create demo $(WASM)
	$(BIN) start demo
	$(BIN) delete demo

bench: build wasm
	@echo "Cold-start x10 (ms):"
	@for i in $$(seq 1 10); do \
		/usr/bin/time -f "%e" $(BIN) create demo$(i) $(WASM) >/dev/null 2>&1; \
		/usr/bin/time -f "%e" $(BIN) start demo$(i) >/dev/null; \
		$(BIN) delete demo$(i) >/dev/null 2>&1; \
	done
