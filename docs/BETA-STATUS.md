# VeGo Beta status — v0.1β+

**Date:** 2026-09-12  
**Line:** optimal BPE-aligned IR (not gzip). Fidelity first; token % measured with tiktoken.

## Compression stack

1. Single-line IR (ASI → `;`)
2. Keyword glyphs = **exactly 1 token** (`cl100k` + `o200k`)
3. Phrase fold (`fmt.Println`, `http.HandleFunc`, `os.Args`, …)
4. **Import-path fold** (`"net/http"` 3→1, `"fmt"` 2→1, …) via codec dictionary
5. Dense glyph juxtaposition (expand re-inserts Go spaces)
6. Σ string table only when a literal repeats (≥2) — avoids preamble loss

## Measured savings (cl100k_base)

| Example | Bytes | Tokens |
|---------|------:|-------:|
| `01_hello` | ~42% | ~16% |
| `02_http` | ~52% | ~27% |
| `03_struct_range` | ~34% | ~13% |
| `04_control` | ~44% | ~25% |

## Still on the table (later)

- Reversible identifier minify when names are multi-token
- Larger learned phrase corpus / domain packs
- ≥60% token milestone as a separate versioned goal
