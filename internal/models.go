package models

// Artist represents the main artist data
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// Relation represents the mapping of dates to locations
type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// IndexData is used for the homepage grid
type IndexData struct {
	Artists []Artist
}

// ArtistFullData aggregates all info for the details page
type ArtistFullData struct {
	Artist   Artist
	Relation Relation
}
