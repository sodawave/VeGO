# VeGo Alpha (pre-release) status

Alpha = pre-versión fidelity cut.

## Done
- Epic 1: keyword glyph IR, expand/compact, round-trip tests, `.vego` fixtures as source artifacts
- Epic 2: `vego fmt|build|run` CLI
- Epic 3: MCP stub with 3 Alpha tools (+ tests)

## Verify
```bash
go test ./src/...
go build -o vego ./src/cmd/vego
./vego fmt src/pkg/transpiler/testdata/hello_http.go > /tmp/x.vego
./vego build /tmp/x.vego
```

## Deferred (post-Alpha / mid-term)
- BPE % goals, tiktoken validation
- Full MCP stdio transport
- AST-true structural patch (substring stand-in now)
- Generics/concurrency grammar
