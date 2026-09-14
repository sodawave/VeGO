#!/usr/bin/env python3
"""Regenerate src/examples/*/stats.md from `vego tokens` (tiktoken) JSON."""
from __future__ import annotations

import json
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]


def main() -> int:
    vego = Path(sys.argv[1]) if len(sys.argv) > 1 else None
    if vego is None:
        bin_path = Path("/tmp/vego-examples-stats")
        subprocess.check_call(
            ["go", "build", "-o", str(bin_path), "./src/cmd/vego"], cwd=ROOT
        )
        vego = bin_path

    examples = ROOT / "src" / "examples"
    for d in sorted(p for p in examples.iterdir() if p.is_dir()):
        go = next(d.glob("*.go"))
        report = json.loads(
            subprocess.check_output([str(vego), "tokens", str(go)], text=True)
        )
        md = f"""# Stats — `{d.name}` (v0.1β)

Measured with **tiktoken `{report.get("estimator", "cl100k_base")}`**.  
Fidelity still required; these % are measured compression, not a contractual SLA.

| Metric | Go | `.vego` | Saving |
|--------|---:|--------:|-------:|
| Bytes | {report["go_bytes"]} | {report["vego_bytes"]} | **{report["byte_saving_pct"]:.1f}%** |
| Tokens (tiktoken) | {report["go_tokens"]} | {report["vego_tokens"]} | **{report["token_saving_pct"]:.1f}%** |

- **Byte ratio** (`.vego` / Go): `{report["byte_ratio"]:.4f}`
- **Token ratio** (`.vego` / Go): `{report["token_ratio"]:.4f}`

## How to regenerate

```bash
python3 src/examples/gen_stats.py
```

## Note

{report.get("note", "")}
"""
        (d / "stats.md").write_text(md)
        print(
            f"{d.name}: bytes {report['byte_saving_pct']:.1f}% | "
            f"tokens {report['token_saving_pct']:.1f}%"
        )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
