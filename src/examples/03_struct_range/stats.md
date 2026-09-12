# Stats — `03_struct_range` (v0.1β)

Measured with **tiktoken `tiktoken/cl100k_base`**.  
Fidelity still required; these % are measured compression, not a contractual SLA.

| Metric | Go | `.vego` | Saving |
|--------|---:|--------:|-------:|
| Bytes | 249 | 156 | **37.3%** |
| Tokens (tiktoken) | 83 | 70 | **15.7%** |

- **Byte ratio** (`.vego` / Go): `0.6265`
- **Token ratio** (`.vego` / Go): `0.8434`

## How to regenerate

```bash
python3 src/examples/gen_stats.py
```

## Note

v0.1β real tiktoken counts. Fidelity still required; token % is measured, not a hard SLA.
