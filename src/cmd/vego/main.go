// Command vego is the CLI entry point for the VeGo toolchain.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sodawave/VeGO/src/pkg/transpiler"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "vego: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		fmt.Println(`vego — VeGo Alpha CLI

Usage:
  vego fmt  <file.vego|file.go>   expand .vego→Go or compact .go→.vego to stdout
  vego build <file.vego> [go build args...]
  vego run   <file.vego> [args...]

Alpha: lossless compact IR; BPE %% is mid-term.`)
		return nil
	}
	cmd, rest := args[0], args[1:]
	switch cmd {
	case "fmt":
		return cmdFmt(rest)
	case "build":
		return cmdBuildRun("build", rest)
	case "run":
		return cmdBuildRun("run", rest)
	default:
		return fmt.Errorf("unknown command %q", cmd)
	}
}

func cmdFmt(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("fmt requires exactly one file")
	}
	path := args[0]
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	eng := transpiler.New(nil)
	var out []byte
	switch {
	case strings.HasSuffix(path, ".vego"):
		out, err = eng.Transform(src, transpiler.ToGo)
	case strings.HasSuffix(path, ".go"):
		out, err = eng.Transform(src, transpiler.ToVeGo)
	default:
		return fmt.Errorf("fmt: unsupported extension (want .vego or .go)")
	}
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(out)
	return err
}

func cmdBuildRun(goCmd string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("%s requires a .vego file", goCmd)
	}
	vegoPath := args[0]
	if !strings.HasSuffix(vegoPath, ".vego") {
		return fmt.Errorf("%s: input must be .vego (got %s)", goCmd, vegoPath)
	}
	src, err := os.ReadFile(vegoPath)
	if err != nil {
		return err
	}
	goSrc, err := transpiler.New(nil).Transform(src, transpiler.ToGo)
	if err != nil {
		return err
	}
	tmpDir, err := os.MkdirTemp("", "vego-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	goFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(goFile, goSrc, 0o644); err != nil {
		return err
	}
	cmdArgs := []string{goCmd, goFile}
	if goCmd == "build" {
		cmdArgs = []string{"build", "-o", strings.TrimSuffix(filepath.Base(vegoPath), ".vego"), goFile}
	}
	if len(args) > 1 {
		if goCmd == "run" {
			cmdArgs = append([]string{"run", goFile, "--"}, args[1:]...)
		} else {
			cmdArgs = append(cmdArgs, args[1:]...)
		}
	}
	c := exec.Command("go", cmdArgs...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Stdin = os.Stdin
	return c.Run()
}
