package main

import (
	"log"
	"net/http"
	"os"

	stock "github.com/ceheng-io/stock-go"
)

func main() {
	addr := os.Getenv("CEHENG_SERVER_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}

	server := NewServer(stock.New())
	log.Printf("ceheng server listening on http://%s", addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatal(err)
	}
}
