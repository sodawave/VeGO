# Stats — `04_control`

Mid-term compression snapshot for this example. **Alpha does not gate on these numbers** (fidelity / round-trip first). Estimator: `vego-rough-v1`.

| Metric | Go | `.vego` | Saving |
|--------|---:|--------:|-------:|
| Bytes | 411 | 282 | **31.4%** |
| Tokens (rough) | 144 | 167 | **-16.0%** |

- **Byte ratio** (`.vego` / Go): `0.6861`
- **Token ratio** (`.vego` / Go): `1.1597`

## How to regenerate

```bash
python3 src/examples/gen_stats.py
# or:
go build -o vego ./src/cmd/vego && ./vego tokens src/examples/04_control/classify.go
```

## Note

`token_saving_pct` uses the Alpha stand-in tokenizer (`vego-rough-v1`): whitespace/ASCII runs + each non-ASCII glyph as one token. Glyphs can **increase** rough token counts even when **bytes drop**. Replace with `tiktoken-go` before treating token % as a product gate.
