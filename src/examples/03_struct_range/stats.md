# Stats — `03_struct_range` (v0.1β)

Measured with **tiktoken `tiktoken/cl100k_base`**.  
Fidelity still required; these % are measured compression, not a contractual SLA.

| Metric | Go | `.vego` | Saving |
|--------|---:|--------:|-------:|
| Bytes | 249 | 174 | **30.1%** |
| Tokens (tiktoken) | 83 | 73 | **12.0%** |

- **Byte ratio** (`.vego` / Go): `0.6988`
- **Token ratio** (`.vego` / Go): `0.8795`

## How to regenerate

```bash
python3 src/examples/gen_stats.py
```

## Note

v0.1β real tiktoken counts. Fidelity still required; token % is measured, not a hard SLA.
