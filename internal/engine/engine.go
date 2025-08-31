package engine

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
)

func RunModule(ctx context.Context, wasmPath string, args []string, env map[string]string, stdin io.Reader, stdout, stderr io.Writer) error {
	wasmBytes, err := os.ReadFile(wasmPath)
	if err != nil {
		return fmt.Errorf("read wasm: %w", err)
	}

	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		return fmt.Errorf("init WASI: %w", err)
	}


	compiled, err := r.CompileModule(ctx, wasmBytes)
	if err != nil {
		return fmt.Errorf("compile: %w", err)
	}
	defer compiled.Close(ctx)

	cfg := wazero.NewModuleConfig().
		WithStdout(stdout).
		WithStderr(stderr).
		WithStdin(stdin).
		WithArgs(args...)

	for k, v := range env {
		cfg = cfg.WithEnv(k, v)
	}

	mod, err := r.InstantiateModule(ctx, compiled, cfg)
	if err != nil {
		// If the module has a start function, instantiate may run it and error here.
		// sys.ExitError captures exit(code) from WASI programs.
		if exit, ok := err.(*sys.ExitError); ok && exit.ExitCode() == 0 {
			return nil
		}
		return fmt.Errorf("instantiate: %w", err)
	}
	defer mod.Close(ctx)

	// Explicitly call _start if exported (common for WASI command modules).
	if start := mod.ExportedFunction("_start"); start != nil {
		if _, err := start.Call(ctx); err != nil {
			if exit, ok := err.(*sys.ExitError); ok && exit.ExitCode() == 0 {
				return nil
			}
			return fmt.Errorf("_start: %w", err)
		}
	}

	return nil
}
