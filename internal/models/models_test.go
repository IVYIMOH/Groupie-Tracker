package models

import (
	"encoding/json"
	"testing"
)

func TestArtistUnmarshaling(t *testing.T) {
	jsonData := `{"id": 1, "name": "Queen", "creationDate": 1970}`
	var artist Artist

	err := json.Unmarshal([]byte(jsonData), &artist)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if artist.Name != "Queen" {
		t.Errorf("Expected Queen, got %s", artist.Name)
	}
}
