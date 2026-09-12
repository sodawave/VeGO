// Command vego is the CLI entry point for the VeGo toolchain.
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "vego: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	_ = args
	fmt.Println("vego: stub — toolchain not implemented yet")
	return nil
}
