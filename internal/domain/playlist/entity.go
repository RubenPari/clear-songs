package playlist

// Playlist represents a playlist entity
type Playlist struct {
	ID       string
	Name     string
	ImageURL string
	Tracks   []Track
}
