# gopache

A caching HTTP reverse proxy written in Go, based on the [Caching Server](https://roadmap.sh/projects/caching-server) project from roadmap.sh.

## Usage

```
gopache --port 3000 --origin https://httpbin.org
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8080` | Port the proxy listens on |
| `--origin` | `""` | Origin server URL to forward requests to |
| `--clear-cache` | `false` | Clear the persistent cache on startup |
| `--cache-path` | `cache.json` | Path to the cache file |

Responses from the origin are cached in memory and persisted to disk (JSON) on graceful shutdown. Subsequent requests for the same resource return the cached response with an `X-Cache: HIT` header, while uncached requests show `X-Cache: MISS`.
