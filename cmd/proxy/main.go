package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func main() {
	handler := proxy.New()

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("🚀 Proxy listening on :8080")
	log.Fatal(server.ListenAndServe())
}
