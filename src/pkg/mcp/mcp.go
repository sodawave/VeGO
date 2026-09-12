// Package mcp stubs the Model Context Protocol tool server for VeGo.
package mcp

import "context"

// ToolRequest is a minimal MCP tool invocation payload.
type ToolRequest struct {
	Name      string
	Arguments map[string]any
}

// ToolResult is a minimal MCP tool response.
type ToolResult struct {
	Content string
}

// Server exposes VeGo read/patch tools to AI agents (stdio or network).
type Server interface {
	Start(ctx context.Context) error
	Handle(ctx context.Context, req ToolRequest) (*ToolResult, error)
}

// Stub is a no-op MCP server until mark3labs/mcp-go wiring lands.
type Stub struct{}

// Start returns nil (clean-slate stub).
func (s *Stub) Start(ctx context.Context) error {
	_ = ctx
	return nil
}

// Handle returns an empty result (clean-slate stub).
func (s *Stub) Handle(ctx context.Context, req ToolRequest) (*ToolResult, error) {
	_ = ctx
	_ = req
	return &ToolResult{Content: ""}, nil
}
