package proxy

import (
	"caching-proxy/internal/cache"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Proxy struct {
	client *http.Client
	cache  cache.Cache
}

func NewProxy(c cache.Cache) *Proxy {
	return &Proxy{
		client: &http.Client{},
		cache:  c,
	}
}

func (p *Proxy) ServeHttp(w http.ResponseWriter, r *http.Request) {
	// Proxy logic to handle requests and utilize cache

	// [1] Determine the target URL.
	originUrl := r.URL.Query().Get("url")
	if originUrl == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}
	// [2] Generate cache key
	cacheKey := r.Method + ":" + originUrl

	// [3] cache lookup (GET only)
	if r.Method == http.MethodGet {

		if data, found := p.cache.Get(cacheKey); found {
			fmt.Println("cached", originUrl)
			w.Write(data)
			return
		}
	}

	// [4] create request to origin server
	req, err := http.NewRequest(r.Method, originUrl, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header = r.Header.Clone()

	// [5] send request to origin server
	resp, err := p.client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	// prevents from open too many files | slow proxy
	defer resp.Body.Close()

	// [6] copy response headers back to client
	for key, value := range resp.Header {
		w.Header()[key] = value
	}

	w.WriteHeader(resp.StatusCode)

	// [7] read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// [8] store in cache (GET only)
	if r.Method == http.MethodGet && resp.StatusCode == http.StatusOK {
		p.cache.Set(cacheKey, body, time.Minute*5) // example TTL of 5 minutes
	}

	// [9] write response body to client
	w.Write(body)
}

// r = the client request
func ProxyHandler(w http.ResponseWriter, r *http.Request) {

	// 2️⃣ Create a new request to the origin server
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// copy headers
	req.Header = r.Header.Clone()

	// 4️⃣ Send request to origin server
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return

	}
	defer resp.Body.Close()

	// 5️⃣ Copy response headers back to the client
	for key, value := range resp.Header {
		w.Header()[key] = value
	}
	// 6️⃣ Send status code + body back
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	//
}
