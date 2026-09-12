# Stats — `02_http` (v0.1β)

Measured with **tiktoken `tiktoken/cl100k_base`**.  
Fidelity still required; these % are measured compression, not a contractual SLA.

| Metric | Go | `.vego` | Saving |
|--------|---:|--------:|-------:|
| Bytes | 240 | 106 | **55.8%** |
| Tokens (tiktoken) | 66 | 48 | **27.3%** |

- **Byte ratio** (`.vego` / Go): `0.4417`
- **Token ratio** (`.vego` / Go): `0.7273`

## How to regenerate

```bash
python3 src/examples/gen_stats.py
```

## Note

v0.1β real tiktoken counts. Fidelity still required; token % is measured, not a hard SLA.
