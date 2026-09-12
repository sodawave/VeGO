// Package mcp implements a minimal MCP-style JSON-RPC tool server for VeGo.
package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	ServeStdio(in io.Reader, out io.Writer) error
}

// Stub is an in-process Alpha MCP handler with JSON-RPC stdio.
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

// Start is a no-op; use ServeStdio for the transport loop.
func (s *Stub) Start(ctx context.Context) error {
	_ = ctx
	return nil
}

// Handle dispatches Alpha tools.
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

type rpcRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
}

type rpcResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   *rpcErr `json:"error,omitempty"`
}

type rpcErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ServeStdio runs a line-delimited JSON-RPC loop (initialize / tools/list / tools/call).
func (s *Stub) ServeStdio(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	enc := json.NewEncoder(out)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var req rpcRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			_ = enc.Encode(rpcResponse{JSONRPC: "2.0", ID: nil, Error: &rpcErr{Code: -32700, Message: err.Error()}})
			continue
		}
		res := s.dispatchRPC(req)
		if err := enc.Encode(res); err != nil {
			return err
		}
	}
	return sc.Err()
}

func (s *Stub) dispatchRPC(req rpcRequest) rpcResponse {
	switch req.Method {
	case "initialize":
		return rpcResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{
			"protocolVersion": "2024-11-05",
			"serverInfo":      map[string]string{"name": "vego-mcp", "version": "alpha"},
			"capabilities":    map[string]any{"tools": map[string]any{}},
		}}
	case "tools/list", "list_tools":
		tools := make([]map[string]any, 0, len(AlphaTools))
		for _, name := range AlphaTools {
			tools = append(tools, map[string]any{
				"name":        name,
				"description": "VeGo Alpha tool " + name,
				"inputSchema": map[string]any{"type": "object"},
			})
		}
		return rpcResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{"tools": tools}}
	case "tools/call", "call_tool":
		name, _ := req.Params["name"].(string)
		args, _ := req.Params["arguments"].(map[string]any)
		if args == nil {
			args = map[string]any{}
		}
		tr, err := s.Handle(context.Background(), ToolRequest{Name: name, Arguments: args})
		if err != nil {
			return rpcResponse{JSONRPC: "2.0", ID: req.ID, Error: &rpcErr{Code: -32000, Message: err.Error()}}
		}
		return rpcResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{
			"content": []map[string]string{{"type": "text", "text": tr.Content}},
		}}
	default:
		return rpcResponse{JSONRPC: "2.0", ID: req.ID, Error: &rpcErr{Code: -32601, Message: "method not found: " + req.Method}}
	}
}
