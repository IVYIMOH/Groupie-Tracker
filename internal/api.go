package handlers

import (
	"encoding/json"
	"net/http"

	"Groupie-Tracker/internal/models"
)

const BaseURL = "https://groupietrackers.herokuapp.com/api"

func FetchArtists() ([]models.Artist, error) {
	resp, err := http.Get(BaseURL + "/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []models.Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	return artists, err
}

func FetchArtistFullData(id string) (models.ArtistDetails, error) {
	var data models.ArtistDetails
	errChan := make(chan error, 2)

	// 1. Fetch Basic Artist Info
	go func() {
		resp, err := http.Get(BaseURL + "/artists/" + id)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()
		// Decodes directly into the embedded Artist struct
		errChan <- json.NewDecoder(resp.Body).Decode(&data.Artist)
	}()

	// 2. Fetch Relations
	go func() {
		resp, err := http.Get(BaseURL + "/relation/" + id)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		// THE FIX: The API returns {"id": 1, "datesLocations": {...}}
		// We need a temporary struct to catch that specific API shape
		var temp struct {
			DatesLocations map[string][]string `json:"datesLocations"`
		}

		err = json.NewDecoder(resp.Body).Decode(&temp)
		data.RelationsMap = temp.DatesLocations
		errChan <- err
	}()

	// Wait for both
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return data, err
		}
	}

	return data, nil
}
