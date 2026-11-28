package main

import (
	"log"
	"net/http"

	"blog-backend/internal/api"
	"blog-backend/internal/store"
)

func main() {
	h := api.NewHandler(store.NewMemory())
	mux := http.NewServeMux()
	h.Register(mux)

	addr := ":8080"
	log.Printf("blog-backend 监听 %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
