# VeGo examples (Alpha)

Runnable Alpha samples under `src/examples/`. Each folder has:

- `*.go` — human-readable Go (expand target / source of truth for humans)
- `*.vego` — compact IR agents emit (generated with `vego fmt`), **single-line** (ASI newlines become `;`)

Alpha grammar only: no generics, no `go`/`chan`/`select`.

## Quick start

```bash
go build -o vego ./src/cmd/vego

# Compact Go → .vego
./vego fmt src/examples/01_hello/hello.go

# Expand .vego → Go (stdout)
./vego fmt src/examples/01_hello/hello.vego

# Build / run from IR
./vego run src/examples/01_hello/hello.vego
./vego build src/examples/01_hello/hello.vego   # writes ./hello in cwd

# Mid-term token estimate (not an Alpha gate)
./vego tokens src/examples/01_hello/hello.go
```

## Catalog

| Dir | What it shows | Run |
|-----|---------------|-----|
| [`01_hello`](./01_hello/) | `package` / `import` / `func` / call | `./vego run src/examples/01_hello/hello.vego` |
| [`02_http`](./02_http/) | HTTP hello (fixture-class) | `./vego run src/examples/02_http/http_server.vego` then open `:8088` (Ctrl+C to stop) |
| [`03_struct_range`](./03_struct_range/) | `type`/`struct` + `for`/`range` | `./vego run src/examples/03_struct_range/sum_points.vego` → prints `10` |
| [`04_control`](./04_control/) | `if`/`else` + counted `for` + args | `./vego run src/examples/04_control/classify.vego 5` → `pos 15` |

## Regenerate `.vego`

After editing a `.go` file:

```bash
./vego fmt src/examples/01_hello/hello.go > src/examples/01_hello/hello.vego
```

Round-trip check: expand the `.vego` and compare semantically with the `.go` (Alpha gate = fidelity).
