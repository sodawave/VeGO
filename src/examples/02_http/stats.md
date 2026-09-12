# Stats — `02_http`

Mid-term compression snapshot for this example. **Alpha does not gate on these numbers** (fidelity / round-trip first). Estimator: `vego-rough-v1`.

| Metric | Go | `.vego` | Saving |
|--------|---:|--------:|-------:|
| Bytes | 240 | 204 | **15.0%** |
| Tokens (rough) | 78 | 87 | **-11.5%** |

- **Byte ratio** (`.vego` / Go): `0.8500`
- **Token ratio** (`.vego` / Go): `1.1154`

## How to regenerate

```bash
go build -o vego ./src/cmd/vego
./vego tokens src/examples/02_http/http_server.go
```

## Note

`token_saving_pct` uses the Alpha stand-in tokenizer (`vego-rough-v1`): whitespace/ASCII runs + each non-ASCII glyph as one token. Glyphs can **increase** rough token counts even when **bytes drop**. Replace with `tiktoken-go` before treating token % as a product gate.

Mid-term heuristic only; Alpha does not gate on token_saving_pct. Swap for tiktoken-go later.
