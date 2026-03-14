package handlers

import (
	"Groupie-Tracker/internal/models" // Adjust this path to match your go.mod
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// FetchArtists fetches the list of artists from the given URL.
func FetchArtists(url string) ([]models.Artist, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET request failed: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var artists []models.Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return artists, nil
}

func FetchArtistDetails(artistID string) (models.ArtistDetails, error) {
	var wg sync.WaitGroup
	var artist models.Artist
	var locations []string
	var dates []string
	var relations map[string][]string
	var err1, err2, err3, err4 error

	wg.Add(4)

	go func() {
		defer wg.Done()
		artist, err1 = FetchArtist(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%s", artistID))
	}()

	go func() {
		defer wg.Done()
		locations, err2 = FetchLocations(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%s", artistID))
	}()

	go func() {
		defer wg.Done()
		dates, err3 = FetchDates(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%s", artistID))
	}()

	go func() {
		defer wg.Done()
		relations, err4 = FetchRelations(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%s", artistID))
	}()

	wg.Wait()

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return models.ArtistDetails{}, fmt.Errorf("one or more fetch operations failed")
	}

	return models.ArtistDetails{
		Artist:           artist,
		LocationsList:    locations,
		ConcertDatesList: dates,
		RelationsMap:     relations,
	}, nil
}

func FetchArtist(url string) (models.Artist, error) {
	response, err := http.Get(url)
	if err != nil {
		return models.Artist{}, err
	}
	defer response.Body.Close()

	var artist models.Artist
	err = json.NewDecoder(response.Body).Decode(&artist)
	return artist, err
}

func FetchLocations(url string) ([]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var data struct {
		Locations []string `json:"locations"`
	}
	err = json.NewDecoder(response.Body).Decode(&data)
	return data.Locations, err
}

func FetchDates(url string) ([]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var data struct {
		Dates []string `json:"dates"`
	}
	err = json.NewDecoder(response.Body).Decode(&data)
	return data.Dates, err
}

func FetchRelations(url string) (map[string][]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var data struct {
		DatesLocations map[string][]string `json:"datesLocations"`
	}
	err = json.NewDecoder(response.Body).Decode(&data)
	return data.DatesLocations, err
}

func SearchArtistsByName(artists []models.Artist, query string) []models.Artist {
	if query == "" {
		return artists
	}
	var filtered []models.Artist
	query = strings.ToLower(query)
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), query) {
			filtered = append(filtered, artist)
		}
	}
	return filtered
}
