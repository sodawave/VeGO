# VeGo Beta status — v0.1β

**Date:** 2026-09-12  
**Gate:** Beta = measured tiktoken compression on the Alpha fidelity base.  
Still **not** a ≥60% SLA; still lossless round-trip first.

## Optimal line shipped

1. **Glyphs = 1 BPE token** in `cl100k_base` and `o200k_base` (Latin-1 / letterlike only).
2. **Phrase fold:** multi-token selectors (`fmt.Println`, `http.HandleFunc`, …) → one glyph.
3. **Single-line IR** (ASI → `;`).
4. **`vego tokens`** reports real **tiktoken** counts (not the rough heuristic).

## Example savings (cl100k_base)

| Example | Byte saving | Token saving |
|---------|------------:|-------------:|
| `01_hello` | ~35% | ~11% |
| `02_http` | ~47% | ~24% |
| `03_struct_range` | ~30% | ~12% |
| `04_control` | ~37% | ~24% |

## Verify

```bash
go test ./src/...
go run ./src/cmd/vego tokens src/examples/02_http/http_server.go
python3 src/examples/gen_stats.py
```

## Still mid-term

- Identifier minify with reversible table
- Broader phrase corpus / learned dictionary
- ≥60% token goal as versioned milestone (not this tag)
- Git textconv / IDE packaging
