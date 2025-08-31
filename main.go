package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wasmrt <command>")
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "create":
		fmt.Println("Creating instance...")
	case "start":
		fmt.Println("Starting instance...")
	default:
		fmt.Println("Unknown command:", cmd)
	}
}
