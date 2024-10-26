package cache

import (
	"time"

	"github.com/RubenPari/clear-songs/src/services/playlistService"

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
		tracks, errGetAllPlaylistTracks := playlistService.GetAllPlaylistTracks(idPlaylist)

		if errGetAllPlaylistTracks != nil {
			return nil, errGetAllPlaylistTracks
		}

		Set("tracksPlaylist"+idPlaylist.String(), tracks)
		playlistTracks = tracks
	}

	return playlistTracks, nil
}
