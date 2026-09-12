# Stats — `03_struct_range`

Mid-term compression snapshot for this example. **Alpha does not gate on these numbers** (fidelity / round-trip first). Estimator: `vego-rough-v1`.

| Metric | Go | `.vego` | Saving |
|--------|---:|--------:|-------:|
| Bytes | 249 | 185 | **25.7%** |
| Tokens (rough) | 87 | 100 | **-14.9%** |

- **Byte ratio** (`.vego` / Go): `0.7430`
- **Token ratio** (`.vego` / Go): `1.1494`

## How to regenerate

```bash
go build -o vego ./src/cmd/vego
./vego tokens src/examples/03_struct_range/sum_points.go
```

## Note

`token_saving_pct` uses the Alpha stand-in tokenizer (`vego-rough-v1`): whitespace/ASCII runs + each non-ASCII glyph as one token. Glyphs can **increase** rough token counts even when **bytes drop**. Replace with `tiktoken-go` before treating token % as a product gate.

Mid-term heuristic only; Alpha does not gate on token_saving_pct. Swap for tiktoken-go later.
