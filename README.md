# logslice

A fast log filtering and time-range slicing utility for large compressed log archives.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build ./...
```

---

## Usage

```bash
logslice [flags] <archive>
```

### Examples

Filter logs within a specific time range:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" /var/logs/app.log.gz
```

Filter by keyword and time range, writing output to a file:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" \
  --filter "ERROR" \
  --out results.log \
  /var/logs/app.log.gz
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--from` | Start of time range (RFC3339) | required |
| `--to` | End of time range (RFC3339) | required |
| `--filter` | Keyword or regex pattern to match | `""` |
| `--out` | Output file path (defaults to stdout) | `-` |
| `--workers` | Number of parallel decompression workers | `4` |

---

## Features

- Supports `.gz`, `.zst`, and `.bz2` compressed archives
- Parallel decompression for faster processing of large files
- Streams output to avoid high memory usage
- Regex and plain-text filtering support

---

## License

MIT © 2024 yourusername