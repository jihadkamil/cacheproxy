package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// r = the client request
// /*
func proxyHandler(w http.ResponseWriter, r *http.Request) {

	// 2️⃣ Create a new request to the origin server
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	fmt.Println("ke sini ga?", err)
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

// */
func testAJa(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("halo ini testing aja")
}

func main() {
	// Client → Proxy → Origin → Proxy → Client
	port := 8000
	// “For every incoming HTTP request, call proxyHandler”
	http.HandleFunc("/", testAJa)
	http.HandleFunc("/prozxy", proxyHandler)

	log.Println("proxy listening on :", port)

	/*
		Starts an HTTP server on port 8080
		Blocks forever
		Every request is routed to proxyHandler
	*/
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
