# Caching Proxy (Go)

A simple HTTP caching proxy written in Go.
This project demonstrates how to build a proxy server that forwards requests to an origin server and caches `GET` responses in memory to reduce repeated upstream calls.

This is a **learning-oriented project** focused on:

* HTTP request/response flow
* Cache proxy fundamentals
* Go interfaces and struct composition

---

## Project Structure

```
cmd/proxy/main.go
internal/cache/cache.go
internal/config/config.go
internal/proxy/proxy.go
```

### Directory Overview

| Path                        | Responsibility                                       |
| --------------------------- | ---------------------------------------------------- |
| `cacheproxy/main.go`         | Application entry point, HTTP server wiring          |
| `cacheproxy/proxy/proxy.go`   | Proxy logic (request forwarding + caching)           |
| `cacheproxy/cache/cache.go`   | Cache interface and in-memory implementation         |
| `cacheproxy/config/config.go` | Reserved for future configuration (currently unused) |

---

## How It Works

```
Client → Proxy → Origin Server
Client ← Proxy ← Origin Server
```

1. The client sends a request to the proxy
2. The proxy extracts the **origin URL** from the `url` query parameter
3. For `GET` requests:

   * the proxy checks the in-memory cache
   * if a cached response exists, it is returned immediately
4. On cache miss:

   * the proxy forwards the request to the origin server
   * reads the response
   * stores it in cache with a TTL
5. The response is returned to the client

Only `GET` requests are cached.
Other HTTP methods are always forwarded to the origin server.

---

## Cache Behavior

* Cache key:

  ```
  <HTTP_METHOD>:<ORIGIN_URL>
  ```
* Cache storage: in-memory map with TTL
* Cache scope: process-local
* Cache eviction: TTL-based only

---

## Example Request

```bash
curl "http://localhost:8000/servehttp?url=http://httpbin.org/get"
```

### First request

* Cache miss
* Request forwarded to origin server
* Response stored in cache

### Second request (same URL)

* Cache hit
* Response returned from cache
* Origin server is not contacted

---

## Running the Project

### Requirements

* Go 1.20 or newer

### Start the proxy

```bash
go run ./cmd/proxy/main.go




```

The proxy listens on port `8000`.

---

## Limitations (By Design)

This project intentionally avoids production-level complexity.

Not implemented:

* Persistent cache (Redis / Memcached)
* Cache invalidation for POST / PUT / DELETE
* Header-aware caching
* Authorization-aware cache keys
* Cache stampede protection (singleflight)
* Metrics and logging

---

## Why This Project Exists

This project is meant to:

* Teach **how caching proxies work**
* Demonstrate clean separation of concerns in Go
* Practice HTTP stream handling and resource safety

It is **not intended for production use** without further hardening.

---

## Possible Next Improvements

If you want to extend this project:

* Add `singleflight` to avoid duplicate concurrent origin requests
* Store response headers alongside the body
* Add configuration (TTL, timeouts) via `internal/config`
* Add Redis-backed cache
* Add structured logging and metrics

---


