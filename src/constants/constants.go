package constants

var Scopes = []string{
	"user-read-private",
	"user-read-email",
	"user-library-read",
	"user-library-modify",
	"playlist-read-private",
	"playlist-read-collaborative",
	"playlist-modify-public",
	"playlist-modify-private",
}

const LimitGetPlaylistTracks = 100
const LimitRemovePlaylistTracks = 100
const LimitInsertPlaylistTracks = 100
const Offset = 0

const PlaylistNameWithMinorSongs = "MinorSongs"
const DescriptionPlaylistNameWithMinorSongs = "A playlist with all tracks from the user library that belong to artists for which the user owns a maximum of 5 tracks."
