package track

// TrackResponse represents a track in API responses
type TrackResponse struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Artists    []string `json:"artists"`
	Album      string   `json:"album"`
	Duration   int      `json:"duration"`
	ImageURL   string   `json:"image_url,omitempty"`
	SpotifyURL string   `json:"spotify_url,omitempty"`
}
