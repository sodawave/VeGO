# VeGo Beta status — v0.1β++

**Date:** 2026-09-12  
**Line:** hard BPE-aligned IR (not gzip). Fidelity first; tiktoken-measured.

## Compression stack

1. Single-line IR (ASI → `;`)
2. Keyword glyphs = **1 token** (`cl100k` + `o200k`)
3. Phrase fold (`fmt.Println`, `http.*`, `os.Args`, …)
4. Import-path fold (`"net/http"` 3→1, …)
5. Dense glyph gluing
6. Σ string table only for **repeated** literals
7. **Hard pass:** `package main`→`¾`, `func main`→`½`, `else if`→`¼`
8. **LitMap:** `"/"` `"--"` `"\n"` `" "`
9. **Assign normalize:** `x = x + y` → `x += y` before encode

## Measured (cl100k_base)

| Example | Bytes | Tokens |
|---------|------:|-------:|
| `01_hello` | ~54% | ~26% |
| `02_http` | ~56% | ~27% |
| `03_struct_range` | ~37% | ~16% |
| `04_control` | ~48% | ~27% |

## Still later

- Project string-bank / domain packs for unique long literals
- Multi-token ident minify with reversible table
- ≥60% token milestone as its own versioned goal
