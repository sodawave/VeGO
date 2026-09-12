package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	var out bytes.Buffer
	if err := run(nil, &out, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "vego fmt") {
		t.Fatalf("help missing fmt: %s", out.String())
	}
}

func TestFmtRoundTripCLI(t *testing.T) {
	goFile := filepath.Join("..", "..", "pkg", "transpiler", "testdata", "hello_http.go")
	src, err := os.ReadFile(goFile)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	inGo := filepath.Join(dir, "hello.go")
	if err := os.WriteFile(inGo, src, 0o644); err != nil {
		t.Fatal(err)
	}
	var vegoOut bytes.Buffer
	if err := run([]string{"fmt", inGo}, &vegoOut, &vegoOut); err != nil {
		t.Fatalf("fmt go→vego: %v", err)
	}
	vegoPath := filepath.Join(dir, "hello.vego")
	if err := os.WriteFile(vegoPath, vegoOut.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	var goOut bytes.Buffer
	if err := run([]string{"fmt", vegoPath}, &goOut, &goOut); err != nil {
		t.Fatalf("fmt vego→go: %v", err)
	}
	if !strings.Contains(goOut.String(), "package main") {
		t.Fatalf("expanded missing package: %s", goOut.String())
	}
}

func TestBuildFromVeGo(t *testing.T) {
	goFile := filepath.Join("..", "..", "pkg", "transpiler", "testdata", "struct_range.go")
	src, err := os.ReadFile(goFile)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	inGo := filepath.Join(dir, "app.go")
	if err := os.WriteFile(inGo, src, 0o644); err != nil {
		t.Fatal(err)
	}
	var vegoOut bytes.Buffer
	if err := run([]string{"fmt", inGo}, &vegoOut, &vegoOut); err != nil {
		t.Fatal(err)
	}
	vegoPath := filepath.Join(dir, "app.vego")
	if err := os.WriteFile(vegoPath, vegoOut.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	cwd, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run([]string{"build", "app.vego"}, &out, &out); err != nil {
		t.Fatalf("build: %v\n%s", err, out.String())
	}
	bin := filepath.Join(dir, "app")
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("missing binary: %v", err)
	}
	cmd := exec.Command(bin)
	if err := cmd.Run(); err != nil {
		t.Fatalf("run binary: %v", err)
	}
}

func TestTokensJSON(t *testing.T) {
	goFile := filepath.Join("..", "..", "pkg", "transpiler", "testdata", "hello_http.go")
	var out bytes.Buffer
	if err := run([]string{"tokens", goFile}, &out, &out); err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(out.Bytes(), &m); err != nil {
		t.Fatalf("json: %v\n%s", err, out.String())
	}
	if m["go_tokens"] == nil || m["vego_tokens"] == nil {
		t.Fatalf("missing fields: %v", m)
	}
}

func TestRunWithArgs(t *testing.T) {
	example := filepath.Join("..", "..", "examples", "04_control", "classify.vego")
	var out bytes.Buffer
	if err := run([]string{"run", example, "5"}, &out, &out); err != nil {
		t.Fatalf("run: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != "pos 15" {
		t.Fatalf("want %q, got %q", "pos 15", got)
	}
}

func TestUnknownCommand(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{"nope"}, &out, &out)
	if err == nil {
		t.Fatal("expected error")
	}
}
