package handlers

import (
	"encoding/json"
	"groupie-tracker/internal/models"
	"net/http"
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

func FetchArtistFullData(id string) (models.ArtistFullData, error) {
	var data models.ArtistFullData
	errChan := make(chan error, 2)

	// Fetch Artist Basic Info
	go func() {
		resp, err := http.Get(BaseURL + "/artists/" + id)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()
		errChan <- json.NewDecoder(resp.Body).Decode(&data.Artist)
	}()

	// Fetch Relations (Locations + Dates)
	go func() {
		resp, err := http.Get(BaseURL + "/relation/" + id)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()
		errChan <- json.NewDecoder(resp.Body).Decode(&data.Relation)
	}()

	// Wait for both routines
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return data, err
		}
	}

	return data, nil
}
