# Stats — `01_hello` (v0.1β)

Measured with **tiktoken `tiktoken/cl100k_base`**.  
Fidelity still required; these % are measured compression, not a contractual SLA.

| Metric | Go | `.vego` | Saving |
|--------|---:|--------:|-------:|
| Bytes | 71 | 41 | **42.3%** |
| Tokens (tiktoken) | 19 | 16 | **15.8%** |

- **Byte ratio** (`.vego` / Go): `0.5775`
- **Token ratio** (`.vego` / Go): `0.8421`

## How to regenerate

```bash
python3 src/examples/gen_stats.py
```

## Note

v0.1β real tiktoken counts. Fidelity still required; token % is measured, not a hard SLA.
