// Command vego is the CLI entry point for the VeGo toolchain.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sodawave/VeGO/src/pkg/bpe"
	"github.com/sodawave/VeGO/src/pkg/mcp"
	"github.com/sodawave/VeGO/src/pkg/transpiler"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "vego: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	_ = stderr
	if len(args) == 0 {
		fmt.Fprintln(stdout, `vego — VeGo Alpha CLI

Usage:
  vego fmt     <file.vego|file.go>   expand .vego→Go or compact .go→.vego
  vego build   <file.vego>           expand then go build
  vego run     <file.vego> [args...] expand then go run
  vego tokens  <file.go|file.vego>   mid-term BPE-ish token estimate (JSON)
  vego mcp     stdio                 JSON-RPC MCP tool server on stdin/stdout

Alpha: lossless compact IR. BPE %% goals are mid-term (tokens command).`)
		return nil
	}
	cmd, rest := args[0], args[1:]
	switch cmd {
	case "fmt":
		return cmdFmt(rest, stdout)
	case "build":
		return cmdBuildRun("build", rest, stdout, stderr)
	case "run":
		return cmdBuildRun("run", rest, stdout, stderr)
	case "tokens":
		return cmdTokens(rest, stdout)
	case "mcp":
		return cmdMCP(rest, os.Stdin, stdout)
	default:
		return fmt.Errorf("unknown command %q", cmd)
	}
}

func cmdFmt(args []string, stdout io.Writer) error {
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
	_, err = stdout.Write(out)
	return err
}

func cmdBuildRun(goCmd string, args []string, stdout, stderr io.Writer) error {
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
	outName := strings.TrimSuffix(filepath.Base(vegoPath), ".vego")
	if goCmd == "build" {
		cmdArgs = []string{"build", "-o", outName, goFile}
	}
	if len(args) > 1 {
		if goCmd == "run" {
			cmdArgs = append([]string{"run", goFile, "--"}, args[1:]...)
		} else {
			cmdArgs = append(cmdArgs, args[1:]...)
		}
	}
	c := exec.Command("go", cmdArgs...)
	c.Stdout = stdout
	c.Stderr = stderr
	c.Stdin = os.Stdin
	return c.Run()
}

func cmdTokens(args []string, stdout io.Writer) error {
	if len(args) != 1 {
		return fmt.Errorf("tokens requires exactly one file")
	}
	path := args[0]
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	eng := transpiler.New(nil)
	var goSrc, vegoSrc []byte
	switch {
	case strings.HasSuffix(path, ".go"):
		goSrc = src
		vegoSrc, err = eng.Transform(src, transpiler.ToVeGo)
	case strings.HasSuffix(path, ".vego"):
		vegoSrc = src
		goSrc, err = eng.Transform(src, transpiler.ToGo)
	default:
		return fmt.Errorf("tokens: want .go or .vego")
	}
	if err != nil {
		return err
	}
	rep := bpe.Compare(goSrc, vegoSrc)
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(rep)
}

func cmdMCP(args []string, in io.Reader, out io.Writer) error {
	if len(args) != 1 || args[0] != "stdio" {
		return fmt.Errorf("usage: vego mcp stdio")
	}
	return mcp.NewStub().ServeStdio(in, out)
}
