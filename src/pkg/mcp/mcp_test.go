package mcp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sodawave/VeGO/src/pkg/mcp"
	"github.com/sodawave/VeGO/src/pkg/transpiler"
)

func TestListTools(t *testing.T) {
	s := mcp.NewStub()
	if len(s.ListTools()) != 3 {
		t.Fatalf("want 3 tools")
	}
}

func TestHandleReadExpandPatch(t *testing.T) {
	dir := t.TempDir()
	goSrc := []byte("package main\n\nfunc main() {}\n")
	vego, err := transpiler.New(nil).Transform(goSrc, transpiler.ToVeGo)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "main.vego")
	if err := os.WriteFile(path, vego, 0o644); err != nil {
		t.Fatal(err)
	}
	s := mcp.NewStub()
	ctx := context.Background()
	r1, err := s.Handle(ctx, mcp.ToolRequest{Name: mcp.ToolReadVeGoContext, Arguments: map[string]any{"path": path}})
	if err != nil || r1.Content == "" {
		t.Fatalf("read: %v", err)
	}
	r2, err := s.Handle(ctx, mcp.ToolRequest{Name: mcp.ToolReadExpandedGo, Arguments: map[string]any{"path": path}})
	if err != nil || !strings.Contains(r2.Content, "package main") {
		t.Fatalf("expand: %v %q", err, r2)
	}
	_, err = s.Handle(ctx, mcp.ToolRequest{
		Name: mcp.ToolStructuralPatch,
		Arguments: map[string]any{"path": path, "find": "ƒ", "replace": "ƒ"},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestServeStdioJSONRPC(t *testing.T) {
	in := strings.NewReader(strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`,
	}, "\n") + "\n")
	var out bytes.Buffer
	if err := mcp.NewStub().ServeStdio(in, &out); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("want 2 responses, got %d: %s", len(lines), out.String())
	}
	var list map[string]any
	if err := json.Unmarshal([]byte(lines[1]), &list); err != nil {
		t.Fatal(err)
	}
	result := list["result"].(map[string]any)
	tools := result["tools"].([]any)
	if len(tools) != 3 {
		t.Fatalf("tools=%v", tools)
	}
}
