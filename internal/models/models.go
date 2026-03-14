package models

// Artist represents the basic information about a music artist.
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

// ArtistDetails extends Artist with additional detailed information.
type ArtistDetails struct {
	Artist
	LocationsList    []string            `json:"locations"`
	ConcertDatesList []string            `json:"concertDates"`
	RelationsMap     map[string][]string `json:"relations"`
}

// PageData holds data for rendering the main artists list page.
type PageData struct {
	Title       string
	Artists     []Artist
	SearchQuery string
}

// ArtistPageData holds data for rendering an individual artist's detail page.
type ArtistPageData struct {
	Title  string
	Artist ArtistDetails
}
