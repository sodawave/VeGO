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

## Catalog (v0.1β+ · tiktoken cl100k_base)

| Dir | What it shows | Byte saving | Token saving | Stats |
|-----|---------------|------------:|-------------:|-------|
| [`01_hello`](./01_hello/) | `package`/`import`/`func` + import-fold + `fmt.Println` | **42.3%** | **15.8%** | [`stats.md`](./01_hello/stats.md) |
| [`02_http`](./02_http/) | HTTP + import/phrase fold | **52.1%** | **27.3%** | [`stats.md`](./02_http/stats.md) |
| [`03_struct_range`](./03_struct_range/) | `type`/`struct` + `for`/`range` | **34.1%** | **13.3%** | [`stats.md`](./03_struct_range/stats.md) |
| [`04_control`](./04_control/) | control flow + `os.Args`/`strconv.Atoi` fold | **44.3%** | **24.8%** | [`stats.md`](./04_control/stats.md) |


## Regenerate `.vego` / `stats.md`

```bash
./vego fmt src/examples/01_hello/hello.go > src/examples/01_hello/hello.vego
python3 src/examples/gen_stats.py
```

Round-trip check: expand the `.vego` and compare semantically with the `.go` (gate = fidelity). Compression % is measured with tiktoken; not a hard ≥60% SLA in v0.1β.
