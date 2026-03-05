package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"groupie-tracker/internal/api"
	"groupie-tracker/internal/models"
)

// Handler holds shared dependencies for all HTTP handlers
type Handler struct {
	client    *api.Client
	templates *template.Template
}

// New creates a Handler, parsing all templates from the templates directory
func New(client *api.Client) (*Handler, error) {
	tmpl, err := template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		return nil, err
	}
	return &Handler{client: client, templates: tmpl}, nil
}

// renderError renders the error page with a given status code and message
func (h *Handler) renderError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	data := struct {
		Code    int
		Message string
	}{Code: status, Message: message}
	if err := h.templates.ExecuteTemplate(w, "error.html", data); err != nil {
		log.Printf("error rendering error page: %v", err)
		http.Error(w, message, status)
	}
}

// Home handles GET / — displays all artists as cards
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.renderError(w, http.StatusNotFound, "Page not found")
		return
	}
	if r.Method != http.MethodGet {
		h.renderError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	artists, err := h.client.GetArtists()
	if err != nil {
		log.Printf("Home: failed to get artists: %v", err)
		h.renderError(w, http.StatusInternalServerError, "Could not load artists. Please try again later.")
		return
	}

	if err := h.templates.ExecuteTemplate(w, "index.html", artists); err != nil {
		log.Printf("Home: template error: %v", err)
	}
}

// ArtistDetail handles GET /artist?id=N — shows full detail for one artist
func (h *Handler) ArtistDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.renderError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		h.renderError(w, http.StatusBadRequest, "Invalid artist ID")
		return
	}

	artist, err := h.client.GetArtist(id)
	if err != nil {
		log.Printf("ArtistDetail: failed to get artist %d: %v", id, err)
		h.renderError(w, http.StatusInternalServerError, "Could not load artist details.")
		return
	}

	relation, err := h.client.GetRelation(id)
	if err != nil {
		log.Printf("ArtistDetail: failed to get relation %d: %v", id, err)
		h.renderError(w, http.StatusInternalServerError, "Could not load concert data.")
		return
	}

	detail := models.ArtistDetail{Artist: artist, Relation: relation}
	if err := h.templates.ExecuteTemplate(w, "artist.html", detail); err != nil {
		log.Printf("ArtistDetail: template error: %v", err)
	}
}

// Search handles GET /search?q=query — returns JSON search results (client-server event)
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	if query == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.SearchResult{})
		return
	}

	artists, err := h.client.GetArtists()
	if err != nil {
		log.Printf("Search: failed to get artists: %v", err)
		http.Error(w, "Search unavailable", http.StatusInternalServerError)
		return
	}

	relations, err := h.client.GetRelations()
	if err != nil {
		log.Printf("Search: failed to get relations: %v", err)
		http.Error(w, "Search unavailable", http.StatusInternalServerError)
		return
	}

	// Build a map of relations by artist ID for quick lookup
	relMap := make(map[int]models.Relation, len(relations.Index))
	for _, rel := range relations.Index {
		relMap[rel.ID] = rel
	}

	var results []models.SearchResult

	for _, a := range artists {
		// Match artist name
		if strings.Contains(strings.ToLower(a.Name), query) {
			results = append(results, models.SearchResult{
				ArtistID: a.ID, Name: a.Name,
				Type: "artist", Value: a.Name,
			})
		}

		// Match members
		for _, m := range a.Members {
			if strings.Contains(strings.ToLower(m), query) {
				results = append(results, models.SearchResult{
					ArtistID: a.ID, Name: a.Name,
					Type: "member", Value: m,
				})
			}
		}

		// Match first album
		if strings.Contains(strings.ToLower(a.FirstAlbum), query) {
			results = append(results, models.SearchResult{
				ArtistID: a.ID, Name: a.Name,
				Type: "firstAlbum", Value: a.FirstAlbum,
			})
		}

		// Match creation date
		if strings.Contains(strconv.Itoa(a.CreationDate), query) {
			results = append(results, models.SearchResult{
				ArtistID: a.ID, Name: a.Name,
				Type: "creationDate", Value: strconv.Itoa(a.CreationDate),
			})
		}

		// Match locations from relation map
		if rel, ok := relMap[a.ID]; ok {
			for loc := range rel.DatesLocations {
				if strings.Contains(strings.ToLower(loc), query) {
					results = append(results, models.SearchResult{
						ArtistID: a.ID, Name: a.Name,
						Type: "location", Value: loc,
					})
					break
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}