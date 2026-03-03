package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"groupie-tracker/internal/models"
)
const baseURL = "https://groupietrackers.herokuapp.com/api"
// Client handles all external API requests
type Client struct {
	httpClient *http.Client
}
//NewClient creates a new API client with a timeout
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},

	}
}

// fetch is a generic helper to GET a URL and decode JSON into dest
func (c *Client) fetch(url string, dest interface{}) error {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("request failed for %s: %w", url, err)

	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)

	}
	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode failed for %s: %w", url, err)

	}
	return nil
} 
//GetArtists returns all artists
func (c *Client) GetArtists() ([]models.Artist, error) {
	var artists []models.Artist
	if err := c.fetch(baseURL+"/artists", &artists); err != nil {
		return nil, err
	}
	return artists, nil
}
//GetArtist returns a single artist by ID
func (c *Client) GetArtist(id int) (models.Artist, error) {
	var artist models.Artist
	if err := c.fetch(fmt.Sprintf("%s/artist/%d", baseURL, id), &artist); err != nil {
		return models.Artist{}, err
	}
	return artist, nil

}
//GetLocations returns all locations
func (c *Client) GetLocations() (models.Locations, error) {
	var locs models.Locations
	if err := c.fetch(baseURL+"/locations", &locs); err != nil {
		return model.Locations{}, err
	}
	return locs, nil
}

//GetDates returns all concert dates
func (c *Client) GetDates() (models.Dates, error) {
	var dates models.Dates
	if err := c.fetch(baseURL+"/dates", &dates); err != nil {
		return models.Dates{}, err

	}
	return dates, nil
}
// GetRelations returns all relation data( links dates+locations to artists)

func (c *Client) GetRelations() (models.Relations, error) {
	var relations models.Relations
	if err := c.fetch(baseURL+"/relation", &relations); err != nil {
		return model.Relations{}, err
	}
	return relations, nil
}

//GetRelations returns the relation for a single artist by ID
func (c *Client) GetRelation(id int) (models.Relation, error) {
	var relation models.Relation
	if err := c.fetch(fmt.Sprintf("%s/relation/%d", baseURL, id), &relation); err != nil {
		return models.Relation{}, err
	}
	return relation, nil
	
}