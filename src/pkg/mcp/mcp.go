// Package mcp stubs the Model Context Protocol tool server for VeGo.
package mcp

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sodawave/VeGO/src/pkg/transpiler"
)

// Tool names locked for Alpha (PRD FR-5).
const (
	ToolReadVeGoContext = "read_vego_context"
	ToolReadExpandedGo  = "read_expanded_go"
	ToolStructuralPatch = "structural_patch"
)

// AlphaTools is the stable Alpha tool list.
var AlphaTools = []string{ToolReadVeGoContext, ToolReadExpandedGo, ToolStructuralPatch}

// ToolRequest is a minimal MCP tool invocation payload.
type ToolRequest struct {
	Name      string
	Arguments map[string]any
}

// ToolResult is a minimal MCP tool response.
type ToolResult struct {
	Content string
}

// Server exposes VeGo read/patch tools to AI agents.
type Server interface {
	Start(ctx context.Context) error
	Handle(ctx context.Context, req ToolRequest) (*ToolResult, error)
	ListTools() []string
}

// Stub is an in-process Alpha MCP handler (stdio transport deferred).
type Stub struct {
	Engine *transpiler.Engine
}

// NewStub returns a Stub with a default transpiler engine.
func NewStub() *Stub {
	return &Stub{Engine: transpiler.New(nil)}
}

// ListTools returns the Alpha tool names.
func (s *Stub) ListTools() []string {
	return append([]string(nil), AlphaTools...)
}

// Start is a no-op until stdio MCP transport is wired.
func (s *Stub) Start(ctx context.Context) error {
	_ = ctx
	return nil
}

// Handle dispatches Alpha tools.
// structural_patch performs a single substring replace (Alpha stand-in for AST patch).
func (s *Stub) Handle(ctx context.Context, req ToolRequest) (*ToolResult, error) {
	_ = ctx
	if s.Engine == nil {
		s.Engine = transpiler.New(nil)
	}
	switch req.Name {
	case ToolReadVeGoContext:
		path, _ := req.Arguments["path"].(string)
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return &ToolResult{Content: string(b)}, nil
	case ToolReadExpandedGo:
		path, _ := req.Arguments["path"].(string)
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		out, err := s.Engine.Transform(b, transpiler.ToGo)
		if err != nil {
			return nil, err
		}
		return &ToolResult{Content: string(out)}, nil
	case ToolStructuralPatch:
		path, _ := req.Arguments["path"].(string)
		find, _ := req.Arguments["find"].(string)
		repl, _ := req.Arguments["replace"].(string)
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if find == "" {
			return nil, fmt.Errorf("structural_patch: find required")
		}
		content := string(b)
		idx := strings.Index(content, find)
		if idx < 0 {
			return nil, fmt.Errorf("structural_patch: find not found")
		}
		content = content[:idx] + repl + content[idx+len(find):]
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return nil, err
		}
		return &ToolResult{Content: "patched"}, nil
	default:
		return nil, fmt.Errorf("unknown tool %q", req.Name)
	}
}
