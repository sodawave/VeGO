# Stats — `04_control` (v0.1β)

Measured with **tiktoken `tiktoken/cl100k_base`**.  
Fidelity still required; these % are measured compression, not a contractual SLA.

| Metric | Go | `.vego` | Saving |
|--------|---:|--------:|-------:|
| Bytes | 411 | 258 | **37.2%** |
| Tokens (tiktoken) | 153 | 117 | **23.5%** |

- **Byte ratio** (`.vego` / Go): `0.6277`
- **Token ratio** (`.vego` / Go): `0.7647`

## How to regenerate

```bash
python3 src/examples/gen_stats.py
```

## Note

v0.1β real tiktoken counts. Fidelity still required; token % is measured, not a hard SLA.
