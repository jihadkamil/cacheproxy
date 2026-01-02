package main

import (
	"caching-proxy/internal/cache"
	"caching-proxy/internal/proxy"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Client → Proxy → Origin → Proxy → Client
	port := 8000
	// “For every incoming HTTP request, call proxyHandler”

	http.HandleFunc("/proxy", proxy.ProxyHandler)
	c := cache.NewInMemoryCache()
	p := proxy.NewProxy(c)
	http.HandleFunc("/servehttp", p.ServeHttp)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
