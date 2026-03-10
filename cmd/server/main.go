package main

import (
	"fmt"
	"groupie-tracker/internal/handlers"
	"net/http"
)

func main() {
	// Serve static files (CSS/JS)
	fs := http.FileServer(http.Dir("internal/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/artist", handlers.ArtistHandler)

	fmt.Println("Server starting at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
