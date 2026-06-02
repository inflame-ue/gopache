# gopache

A caching HTTP reverse proxy written in Go, based on the [Caching Server](https://roadmap.sh/projects/caching-server) project from roadmap.sh.

## Usage

First, build the project with

```bash
go build -o gopache ./cmd/proxy
```

Then you can run the proxy with

```
./gopache --port 3000 --origin https://httpbin.org
```

or make it avalaible in `$PATH` for even easier use. 

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8080` | Port the proxy listens on |
| `--origin` | `""` | Origin server URL to forward requests to |
| `--clear-cache` | `false` | Clear the persistent cache on startup |
| `--cache-path` | `cache.json` | Path to the cache file |

Responses from the origin are cached in memory and persisted to disk (JSON) on graceful shutdown. Subsequent requests for the same resource return the cached response with an `X-Cache: HIT` header, while uncached requests show `X-Cache: MISS`.
