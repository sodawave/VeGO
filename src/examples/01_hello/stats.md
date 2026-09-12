# Stats — `01_hello`

Mid-term compression snapshot for this example. **Alpha does not gate on these numbers** (fidelity / round-trip first). Estimator: `vego-rough-v1`.

| Metric | Go | `.vego` | Saving |
|--------|---:|--------:|-------:|
| Bytes | 71 | 55 | **22.5%** |
| Tokens (rough) | 21 | 25 | **-19.0%** |

- **Byte ratio** (`.vego` / Go): `0.7746`
- **Token ratio** (`.vego` / Go): `1.1905`

## How to regenerate

```bash
go build -o vego ./src/cmd/vego
./vego tokens src/examples/01_hello/hello.go
```

## Note

`token_saving_pct` uses the Alpha stand-in tokenizer (`vego-rough-v1`): whitespace/ASCII runs + each non-ASCII glyph as one token. Glyphs can **increase** rough token counts even when **bytes drop**. Replace with `tiktoken-go` before treating token % as a product gate.

Mid-term heuristic only; Alpha does not gate on token_saving_pct. Swap for tiktoken-go later.
