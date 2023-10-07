package models

type ArtistGroupSummary struct {
	Genre   string          `json:"genre"`
	Artists []ArtistSummary `json:"artists"`
}
