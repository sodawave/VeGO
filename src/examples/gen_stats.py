#!/usr/bin/env python3
"""Regenerate src/examples/*/stats.md from `vego tokens` JSON."""
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
        subprocess.check_call(["go", "build", "-o", str(bin_path), "./src/cmd/vego"], cwd=ROOT)
        vego = bin_path

    examples = ROOT / "src" / "examples"
    for d in sorted(p for p in examples.iterdir() if p.is_dir()):
        go = next(d.glob("*.go"))
        report = json.loads(
            subprocess.check_output([str(vego), "tokens", str(go)], text=True)
        )
        byte_save = (1 - report["byte_ratio"]) * 100
        tok_save = report["token_saving_pct"]
        md = f"""# Stats — `{d.name}`

Mid-term compression snapshot for this example. **Alpha does not gate on these numbers** (fidelity / round-trip first). Estimator: `{report["estimator"]}`.

| Metric | Go | `.vego` | Saving |
|--------|---:|--------:|-------:|
| Bytes | {report["go_bytes"]} | {report["vego_bytes"]} | **{byte_save:.1f}%** |
| Tokens (rough) | {report["go_tokens"]} | {report["vego_tokens"]} | **{tok_save:.1f}%** |

- **Byte ratio** (`.vego` / Go): `{report["byte_ratio"]:.4f}`
- **Token ratio** (`.vego` / Go): `{report["token_ratio"]:.4f}`

## How to regenerate

```bash
python3 src/examples/gen_stats.py
# or:
go build -o vego ./src/cmd/vego && ./vego tokens {go.relative_to(ROOT).as_posix()}
```

## Note

`token_saving_pct` uses the Alpha stand-in tokenizer (`{report["estimator"]}`): whitespace/ASCII runs + each non-ASCII glyph as one token. Glyphs can **increase** rough token counts even when **bytes drop**. Replace with `tiktoken-go` before treating token % as a product gate.
"""
        (d / "stats.md").write_text(md)
        print(f"{d.name}: bytes {byte_save:.1f}% | tokens {tok_save:.1f}%")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
