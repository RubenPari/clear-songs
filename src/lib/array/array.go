package array

import (
	"github.com/RubenPari/clear-songs/src/models"
)

// Contains checks if an array of string
// contains an element string
func Contains(array []string, element string) bool {
	for _, a := range array {
		if a == element {
			return true
		}
	}

	return false
}

// ContainsGenre checks if an array of string genres
// contains almost one genre of the second array of genres
func ContainsGenre(genres []string, genresToSearch []string) bool {
	for _, genre := range genres {
		if Contains(genresToSearch, genre) {
			return true
		}
	}

	return false
}

// FilterByMin returns an array of tracks
// of artist that have at least the
// minimum number of tracks
func FilterByMin(tracks map[string]int, min int) map[string]int {
	var newTracks = make(map[string]int)

	for artist, count := range tracks {
		if count >= min {
			newTracks[artist] = count
		}
	}

	return newTracks
}

// FilterByMax returns an array of tracks
// of artist that have at most the
// maximum number of tracks
func FilterByMax(tracks map[string]int, max int) map[string]int {
	var newTracks = make(map[string]int)

	for artist, count := range tracks {
		if count <= max {
			newTracks[artist] = count
		}
	}

	return newTracks
}

// FilterSummaryByRange returns an array of
// artist summary that have at least the
// minimum number of tracks and at most the
// maximum number of tracks
// NOTE: if min or max are 0, they are ignored
func FilterSummaryByRange(tracks []models.ArtistSummary, min int, max int) []models.ArtistSummary {
	var newTracks []models.ArtistSummary

	for _, track := range tracks {
		if min == 0 && max == 0 {
			newTracks = append(newTracks, track)
		} else if min == 0 && track.Count <= max {
			newTracks = append(newTracks, track)
		} else if max == 0 && track.Count >= min {
			newTracks = append(newTracks, track)
		} else if track.Count >= min && track.Count <= max {
			newTracks = append(newTracks, track)
		}
	}

	return newTracks
}
