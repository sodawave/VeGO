package mcp_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/sodawave/VeGO/src/pkg/mcp"
	"github.com/sodawave/VeGO/src/pkg/transpiler"
)

func TestListTools(t *testing.T) {
	s := mcp.NewStub()
	tools := s.ListTools()
	if len(tools) != 3 {
		t.Fatalf("want 3 tools, got %v", tools)
	}
}

func TestReadAndExpand(t *testing.T) {
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
		t.Fatalf("read_vego_context: %v %#v", err, r1)
	}
	r2, err := s.Handle(ctx, mcp.ToolRequest{Name: mcp.ToolReadExpandedGo, Arguments: map[string]any{"path": path}})
	if err != nil {
		t.Fatal(err)
	}
	if !contains(r2.Content, "package main") {
		t.Fatalf("expanded missing package: %q", r2.Content)
	}
	_, err = s.Handle(ctx, mcp.ToolRequest{
		Name: mcp.ToolStructuralPatch,
		Arguments: map[string]any{
			"path":    path,
			"find":    "ƒ",
			"replace": "ƒ",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || stringIndex(s, sub) >= 0)
}

func stringIndex(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
