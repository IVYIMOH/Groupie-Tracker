package handlers

import (
	"groupie-tracker/internal/models"
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	artists, err := FetchArtists()
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("internal/templates/index.html"))
	tmpl.Execute(w, models.IndexData{Artists: artists})
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	fullData, err := FetchArtistFullData(id)
	if err != nil {
		http.Error(w, "Artist Not Found", http.StatusNotFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("internal/templates/artist.html"))
	tmpl.Execute(w, fullData)
}
