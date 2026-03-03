package main

import (
	"log"
	"net/http"

	"groupie-tracker/internal/api"
	"groupie-tracker/internal/handlers"
)

func main() {
	client := api.NewClient()

	h, err := handlers.New(client)
	if err != nil {
		log.Fatalf("Failed to load templates: %v", err)

	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.Home)
	mux.HandleFunc("/artist", h.ArtistDetail)
	mux.HandleFunc("/search", h.Search)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}