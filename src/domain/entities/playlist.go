package entities

// Playlist represents a playlist entity
type Playlist struct {
	ID       string
	Name     string
	ImageURL string
	Tracks   []Track
}
