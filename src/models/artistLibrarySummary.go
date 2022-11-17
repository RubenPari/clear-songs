package models

// ArtistLibrarySummary is a model that contains
// the following information:
// - id artist that has almost 1 song in the user library
// - name of the artist defined above
// - number of songs that the artist has in the user library
type ArtistLibrarySummary struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Num  int    `json:"num"`
}
