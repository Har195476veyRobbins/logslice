# logslice

A fast log filtering and aggregation utility with support for structured JSON and plain-text log formats.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build -o logslice .
```

---

## Usage

Filter logs by level and time range:

```bash
logslice --level error --since 1h ./app.log
```

Aggregate JSON log fields and output a summary:

```bash
logslice --format json --group-by service --count ./logs/*.log
```

Pipe from stdin:

```bash
cat app.log | logslice --level warn --format plain
```

### Flags

| Flag | Description |
|-------------|--------------------------------------|
| `--level` | Filter by log level (info, warn, error) |
| `--since` | Show logs from the last duration (e.g. `30m`, `2h`) |
| `--format` | Input format: `json` or `plain` (default: `plain`) |
| `--group-by` | Aggregate results by a JSON field key |
| `--count` | Print occurrence counts per group |

---

## License

[MIT](LICENSE)