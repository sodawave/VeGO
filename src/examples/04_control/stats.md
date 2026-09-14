# Stats — `04_control` (v0.1β)

Measured with **tiktoken `tiktoken/cl100k_base`**.  
Fidelity still required; these % are measured compression, not a contractual SLA.

| Metric | Go | `.vego` | Saving |
|--------|---:|--------:|-------:|
| Bytes | 411 | 213 | **48.2%** |
| Tokens (tiktoken) | 153 | 111 | **27.5%** |

- **Byte ratio** (`.vego` / Go): `0.5182`
- **Token ratio** (`.vego` / Go): `0.7255`

## How to regenerate

```bash
python3 src/examples/gen_stats.py
```

## Note

v0.1β real tiktoken counts. Fidelity still required; token % is measured, not a hard SLA.
