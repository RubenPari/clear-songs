package array

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

// FilterSummaryByRange returns a
// map with the number of tracks
// of each artist that have at least
// the minimum number and maximum
func FilterSummaryByRange(tracks map[string]int, min int, max int) map[string]int {
	var newSummary = make(map[string]int)

	for artist, count := range tracks {
		if count >= min && count <= max {
			newSummary[artist] = count
		}
	}

	return newSummary
}
