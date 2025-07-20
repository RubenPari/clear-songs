package cache

import (
	"log"
	"time"

	"github.com/RubenPari/clear-songs/src/services/playlistService"
	"github.com/RubenPari/clear-songs/src/services/userService"

	"github.com/patrickmn/go-cache"
	spotifyAPI "github.com/zmb3/spotify"
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
func Get(key string) interface{} {
	value, found := cacheStore.Get(key)
	if found {
		return value
	}
	return nil
}

// InvalidateUserData invalidates all user-related cache entries
// This should be called after any operation that modifies user's tracks
func InvalidateUserData() {
	log.Println("Invalidating user data cache")
	cacheStore.Delete("userTracks")

	// Also invalidate any playlist cache since user modifications
	// might affect playlist content
	InvalidateAllPlaylists()
}

// InvalidatePlaylist invalidates cache for a specific playlist
func InvalidatePlaylist(playlistID spotifyAPI.ID) {
	log.Printf("Invalidating cache for playlist: %s", playlistID.String())
	cacheStore.Delete("tracksPlaylist" + playlistID.String())
}

// InvalidateAllPlaylists invalidates all playlist caches
func InvalidateAllPlaylists() {
	log.Println("Invalidating all playlist caches")

	// Get all cache items and delete playlist-related ones
	items := cacheStore.Items()
	for key := range items {
		if len(key) > 14 && key[:14] == "tracksPlaylist" {
			cacheStore.Delete(key)
		}
	}
}

// Reset clears all cache
func Reset() {
	log.Println("Resetting entire cache")
	cacheStore.Flush()
}

// GetCachedPlaylistTracksOrSet retrieves a list of tracks from the cache, or if not present,
// fetches the tracks from the Spotify API and stores them in the cache.
func GetCachedPlaylistTracksOrSet(idPlaylist spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	var playlistTracks []spotifyAPI.PlaylistTrack

	cacheKey := "tracksPlaylist" + idPlaylist.String()
	value := Get(cacheKey)

	if value != nil {
		log.Printf("Cache hit for playlist: %s", idPlaylist.String())
		playlistTracks = value.([]spotifyAPI.PlaylistTrack)
	} else {
		log.Printf("Cache miss for playlist: %s, fetching from API", idPlaylist.String())
		tracks, errGetAllPlaylistTracks := playlistService.GetAllPlaylistTracks(idPlaylist)

		if errGetAllPlaylistTracks != nil {
			return nil, errGetAllPlaylistTracks
		}

		Set(cacheKey, tracks)
		playlistTracks = tracks
	}

	return playlistTracks, nil
}

// GetCachedUserTracksOrSet retrieves user tracks from cache or fetches them
func GetCachedUserTracksOrSet() ([]spotifyAPI.SavedTrack, error) {
	var userTracks []spotifyAPI.SavedTrack

	value := Get("userTracks")

	if value != nil {
		log.Println("Cache hit for user tracks")
		userTracks = value.([]spotifyAPI.SavedTrack)
	} else {
		log.Println("Cache miss for user tracks, fetching from API")
		tracks, errTracks := userService.GetAllUserTracks()

		if errTracks != nil {
			return nil, errTracks
		}

		Set("userTracks", tracks)
		userTracks = tracks
	}

	return userTracks, nil
}
