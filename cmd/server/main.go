package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"Groupie-Tracker/internal/handlers"
	"Groupie-Tracker/internal/models"
)

// Package-level variable to hold the data fetched at startup
var artistsData []models.Artist

// --- HANDLERS ---

func renderStatusPage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	statusTemplate, err := template.ParseFiles("templates/status.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	statusTemplate.Execute(w, message)
}

// HomeHandler is now a named function, making it accessible to tests.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderStatusPage(w, http.StatusNotFound, "404 Page Not Found")
		return
	}

	query := r.URL.Query().Get("search")
	filteredArtists := handlers.SearchArtistsByName(artistsData, query)

	data := models.PageData{
		Title:       "Groupie-Tracker",
		Artists:     filteredArtists,
		SearchQuery: query,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		renderStatusPage(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	tmpl.Execute(w, data)
}

// ArtistHandler handles specific artist details.
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	artistID := r.URL.Path[len("/artist/"):]
	artist, err := handlers.FetchArtistDetails(artistID)
	if err != nil {
		renderStatusPage(w, http.StatusBadRequest, "400 Bad Request")
		return
	}

	data := models.ArtistPageData{
		Title:  artist.Name,
		Artist: artist,
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		renderStatusPage(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	tmpl.Execute(w, data)
}

// --- FILE SERVERS ---

func safeStaticHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(dir, r.URL.Path[len("/static/"):])
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			renderStatusPage(w, http.StatusNotFound, "404 Page Not Found")
			return
		}
		http.ServeFile(w, r, path)
	}
}

// --- MAIN ---

func main() {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	var err error

	// Initial data load
	artistsData, err = handlers.FetchArtists(url)
	if err != nil {
		log.Fatalf("Critical Error: Failed to fetch artists: %s", err)
	}

	http.HandleFunc("/static/", safeStaticHandler("static"))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/artist/", ArtistHandler)

	log.Println("\033[32mWelcome To Groupie-Tracker\033[0m\nServer started on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
