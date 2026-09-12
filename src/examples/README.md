# VeGo examples (v0.1β)

Runnable Beta samples under `src/examples/`. Each folder has:

- `*.go` — human-readable Go (expand target / source of truth for humans)
- `*.vego` — compact IR (single-line; 1-token glyphs + phrase fold)
- `stats.md` — **tiktoken cl100k_base** byte/token saving %

Grammar subset: no generics, no `go`/`chan`/`select`.

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

# Mid-term token estimate (tiktoken, Beta)
./vego tokens src/examples/01_hello/hello.go
```

## Catalog (v0.1β · tiktoken cl100k_base)

| Dir | What it shows | Byte saving | Token saving | Stats |
|-----|---------------|------------:|-------------:|-------|
| [`01_hello`](./01_hello/) | `package` / `import` / `func` / call + `fmt.Println` phrase | **35.2%** | **10.5%** | [`stats.md`](./01_hello/stats.md) |
| [`02_http`](./02_http/) | HTTP hello + phrase fold (`HandleFunc`, `ListenAndServe`, …) | **46.7%** | **24.2%** | [`stats.md`](./02_http/stats.md) |
| [`03_struct_range`](./03_struct_range/) | `type`/`struct` + `for`/`range` | **30.1%** | **12.0%** | [`stats.md`](./03_struct_range/stats.md) |
| [`04_control`](./04_control/) | `if`/`else` + counted `for` + args + `strconv.Atoi` | **37.2%** | **23.5%** | [`stats.md`](./04_control/stats.md) |


## Regenerate `.vego` / `stats.md`

```bash
./vego fmt src/examples/01_hello/hello.go > src/examples/01_hello/hello.vego
python3 src/examples/gen_stats.py
```

Round-trip check: expand the `.vego` and compare semantically with the `.go` (gate = fidelity). Compression % is measured with tiktoken; not a hard ≥60% SLA in v0.1β.
