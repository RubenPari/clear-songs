package cache

import (
	"time"

	"github.com/RubenPari/clear-songs/src/services"
	spotifyAPI "github.com/zmb3/spotify"

	"github.com/patrickmn/go-cache"
)

var cacheStore *cache.Cache

// Init creates a cacheManager with a default
// expiration time of 5 and which purges
// expired items every 10
func Init() {
	cacheStore = cache.New(5*time.Minute, 10*time.Minute)
}

// Set adds a new item to the cache
// The item is added with the default expiration time
func Set(key string, value interface{}) {
	cacheStore.Set(key, value, cache.DefaultExpiration)
}

// Get retrieves an item from the cache
// If the item does not exist, it returns nil
func Get(key string) (interface{}, bool) {
	value, found := cacheStore.Get(key)
	if found {
		return value, true
	}
	return nil, false
}

// GetCachedAlbumsOrSet retrieves a list of albums from the cache, or if not present,
// fetches the albums from the Spotify API and stores them in the cache.
// It returns the list of albums.
func GetCachedAlbumsOrSet() []spotifyAPI.SavedAlbum {
	var albums []spotifyAPI.SavedAlbum

	value, found := Get("albums")

	if found {
		albums = value.([]spotifyAPI.SavedAlbum)
	} else {
		albums = services.GetAllUserAlbums()

		Set("albums", albums)
	}

	return albums
}

// GetCachedPlaylistTracksOrSet retrieves a list of tracks from the cache, or if not present,
// fetches the tracks from the Spotify API and stores them in the cache.
// It returns the list of tracks.
// The cache is stored with the key "tracksPlaylist" + idPlaylist.
func GetCachedPlaylistTracksOrSet(idPlaylist spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	var playlistTracks []spotifyAPI.PlaylistTrack

	value, found := Get("tracksPlaylist" + idPlaylist.String())

	if found {
		playlistTracks = value.([]spotifyAPI.PlaylistTrack)
	} else {
		tracks, errGetAllPlaylistTracks := services.GetAllPlaylistTracks(idPlaylist)

		if errGetAllPlaylistTracks != nil {
			return nil, errGetAllPlaylistTracks
		}

		Set("tracksPlaylist"+idPlaylist.String(), tracks)
		playlistTracks = tracks
	}

	return playlistTracks, nil
}
