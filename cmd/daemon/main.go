package main

import (
	"log"
	ht "net/http"

	"github.com/myrat012/testing-websocket-golang/internal/controller/http"
)

func main() {
	r := http.NewRouter()
	ht.Handle("/", r)
	log.Println("Starting server on :8080")
	log.Fatal(ht.ListenAndServe(":8080", nil))
}
