package main

import (
	"Groupie-Tracker/internal/handlers" // For the fetch functions
	"Groupie-Tracker/internal/models"   // For the data structures
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// renderStatusPage stays the same, just ensures correct template path
func renderStatusPage(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	// Make sure your templates folder is in the root directory!
	statusTemplate, err := template.ParseFiles("templates/status.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	statusTemplate.Execute(w, message)
}

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

func safeImageHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(dir, r.URL.Path[len("/images/"):])
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			renderStatusPage(w, http.StatusNotFound, "404 Page Not Found")
			return
		}
		http.ServeFile(w, r, path)
	}
}

func main() {
	url := "https://groupietrackers.herokuapp.com/api/artists"

	// Fetch artists using the handlers package
	artists, err := handlers.FetchArtists(url)
	if err != nil {
		log.Fatalf("Failed to fetch artists: %s", err)
	}

	http.HandleFunc("/static/", safeStaticHandler("static"))
	http.HandleFunc("/images/", safeImageHandler("images"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			renderStatusPage(w, http.StatusNotFound, "404 Page Not Found")
			return
		}

		query := r.URL.Query().Get("search")
		filteredArtists := handlers.SearchArtistsByName(artists, query)

		// FIX: Use models.PageData
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
	})

	http.HandleFunc("/artist/", func(w http.ResponseWriter, r *http.Request) {
		artistID := r.URL.Path[len("/artist/"):]

		// Fetch details using the handlers package
		artist, err := handlers.FetchArtistDetails(artistID)
		if err != nil {
			renderStatusPage(w, http.StatusBadRequest, "400 Bad Request")
			return
		}

		// FIX: Use models.ArtistPageData
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
	})

	log.Println("\033[32mWelcome To Groupie-Tracker\033[0m\nServer started on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
