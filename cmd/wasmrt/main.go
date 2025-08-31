package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/deepankar20/wasm-runtime-BTP/internal/engine"
)

type Container struct {
	ID      string            `json:"id"`
	Module  string            `json:"module"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
	Created time.Time         `json:"created"`
}

const stateDir = ".state"

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	switch os.Args[1] {
	case "create":
		cmdCreate(os.Args[2:])
	case "start":
		cmdStart(os.Args[2:])
	case "delete":
		cmdDelete(os.Args[2:])
	default:
		usage()
	}
}

func usage() {
	fmt.Println(`Usage:
  wasmrt create <id> <module.wasm> [--arg value ...] [--env KEY=VAL ...]
  wasmrt start  <id> [--timeout 2s]
  wasmrt delete <id>
`)
}

func cmdCreate(args []string) {
	if len(args) < 2 {
		fmt.Println("create requires <id> <module.wasm>")
		os.Exit(1)
	}
	id := args[0]
	mod := args[1]
	rest := args[2:]

	c := &Container{
		ID:      id,
		Module:  mod,
		Args:    []string{},
		Env:     map[string]string{},
		Created: time.Now(),
	}

	for i := 0; i < len(rest); i++ {
		switch rest[i] {
		case "--arg":
			if i+1 >= len(rest) {
				fmt.Println("--arg needs a value")
				os.Exit(1)
			}
			c.Args = append(c.Args, rest[i+1])
			i++
		case "--env":
			if i+1 >= len(rest) {
				fmt.Println("--env needs KEY=VAL")
				os.Exit(1)
			}
			k, v, ok := splitKV(rest[i+1])
			if !ok {
				fmt.Println("invalid --env, expected KEY=VAL")
				os.Exit(1)
			}
			c.Env[k] = v
			i++
		default:
			fmt.Println("unknown option:", rest[i])
			os.Exit(1)
		}
	}

	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		panic(err)
	}
	data, _ := json.MarshalIndent(c, "", "  ")
	if err := os.WriteFile(filepath.Join(stateDir, id+".json"), data, 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("created %q -> %s\n", id, c.Module)
}

func cmdStart(args []string) {
	if len(args) < 1 {
		fmt.Println("start requires <id>")
		os.Exit(1)
	}
	id := args[0]
	timeout := 0 * time.Second
	rest := args[1:]
	for i := 0; i < len(rest); i++ {
		switch rest[i] {
		case "--timeout":
			if i+1 >= len(rest) {
				fmt.Println("--timeout needs a duration, e.g. 2s")
				os.Exit(1)
			}
			d, err := time.ParseDuration(rest[i+1])
			if err != nil {
				fmt.Println("invalid duration:", err)
				os.Exit(1)
			}
			timeout = d
			i++
		default:
			fmt.Println("unknown option:", rest[i])
			os.Exit(1)
		}
	}

	path := filepath.Join(stateDir, id+".json")
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("not found:", id)
		os.Exit(1)
	}
	var c Container
	if err := json.Unmarshal(b, &c); err != nil {
		panic(err)
	}

	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	if err := engine.RunModule(ctx, c.Module, c.Args, c.Env, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "module error:", err)
		os.Exit(1)
	}
	fmt.Printf("started %q\n", id)
}

func cmdDelete(args []string) {
	if len(args) < 1 {
		fmt.Println("delete requires <id>")
		os.Exit(1)
	}
	id := args[0]
	path := filepath.Join(stateDir, id+".json")
	if err := os.Remove(path); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("deleted %q\n", id)
}

func splitKV(s string) (string, string, bool) {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return s[:i], s[i+1:], true
		}
	}
	return "", "", false
}
