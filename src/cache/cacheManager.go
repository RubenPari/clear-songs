package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/RubenPari/clear-songs/src/services/playlistService"
	"github.com/RubenPari/clear-songs/src/services/userService"

	"github.com/redis/go-redis/v9"
	spotifyAPI "github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var (
	rdb         *redis.Client
	ctx         = context.Background()
	memoryToken *oauth2.Token // Fallback in-memory token storage when Redis is not available
)

const (
	defaultTTL = 5 * time.Minute
	tokenTTL   = 24 * time.Hour // Token should last 24 hours
)

// Init initializes a Redis client using environment variables
// REDIS_HOST, REDIS_PORT, REDIS_PASSWORD, REDIS_DB
func Init() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}
	db := 0
	if dbStr != "" {
		fmt.Sscanf(dbStr, "%d", &db)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	// Try to ping Redis, but don't fail if it's not available
	// This allows the application to run without Redis (though with reduced performance)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("WARNING: Redis connection failed: %v", err)
		log.Println("WARNING: Application will continue without Redis caching. Some features may be slower.")
		log.Println("WARNING: To enable Redis, start Redis server and restart the application.")
		rdb = nil // Set to nil to indicate Redis is not available
		return
	}

	log.Println("Connected to Redis for caching")
}

// internal helpers
func setJSON(key string, value interface{}, ttl time.Duration) error {
	if rdb == nil {
		return nil // Redis not available, silently fail (no caching)
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, data, ttl).Err()
}

func getJSON(key string, target interface{}) (bool, error) {
	if rdb == nil {
		return false, nil // Redis not available, return not found
	}
	val, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	if err := json.Unmarshal(val, target); err != nil {
		return false, err
	}
	return true, nil
}

// Token cache
func SetToken(token *oauth2.Token) error {
	// If token is nil, clear it
	if token == nil {
		ClearToken()
		return nil
	}
	
	if rdb == nil {
		// Redis not available, use in-memory fallback
		log.Println("WARNING: Redis not available, using in-memory token storage (not persistent across restarts)")
		memoryToken = token
		return nil
	}
	err := setJSON("spotify_token", token, tokenTTL)
	if err != nil {
		log.Printf("ERROR: Failed to save token to Redis: %v", err)
		// Fallback to in-memory storage
		memoryToken = token
		return nil
	}
	// Clear in-memory token when Redis is working
	memoryToken = nil
	log.Println("Token saved to Redis cache successfully")
	return nil
}

func GetToken() *oauth2.Token {
	// First try Redis if available
	if rdb != nil {
		var token oauth2.Token
		ok, err := getJSON("spotify_token", &token)
		if err != nil {
			log.Printf("ERROR: Failed to retrieve token from Redis: %v", err)
			// Fallback to in-memory token
			return memoryToken
		}
		if ok {
			log.Println("Token retrieved from Redis cache successfully")
			return &token
		}
	}
	
	// Redis not available or token not found in Redis, try in-memory fallback
	if memoryToken != nil {
		log.Println("Token retrieved from in-memory storage")
		return memoryToken
	}
	
	// No token found
	return nil
}

// ClearToken clears the token from both Redis and in-memory storage
func ClearToken() {
	if rdb != nil {
		_ = rdb.Del(ctx, "spotify_token").Err()
	}
	memoryToken = nil
	log.Println("Token cleared from cache")
}

// InvalidateUserData invalidates all user-related cache entries
func InvalidateUserData() {
	if rdb == nil {
		return // Redis not available
	}
	log.Println("Invalidating user data cache")
	_ = rdb.Del(ctx, "userTracks").Err()
	InvalidateAllPlaylists()
}

// InvalidatePlaylist invalidates cache for a specific playlist
func InvalidatePlaylist(playlistID spotifyAPI.ID) {
	if rdb == nil {
		return // Redis not available
	}
	log.Printf("Invalidating cache for playlist: %s", playlistID.String())
	_ = rdb.Del(ctx, "tracksPlaylist"+playlistID.String()).Err()
}

// InvalidateAllPlaylists invalidates all playlist caches
func InvalidateAllPlaylists() {
	if rdb == nil {
		return // Redis not available
	}
	log.Println("Invalidating all playlist caches")
	iter := rdb.Scan(ctx, 0, "tracksPlaylist*", 0).Iterator()
	for iter.Next(ctx) {
		_ = rdb.Del(ctx, iter.Val()).Err()
	}
}

// Reset clears all cache
func Reset() {
	log.Println("Resetting entire cache")
	// Best-effort: delete known prefixes
	InvalidateUserData()
}

// GetCachedPlaylistTracksOrSet retrieves a list of tracks from cache or API
func GetCachedPlaylistTracksOrSet(idPlaylist spotifyAPI.ID) ([]spotifyAPI.PlaylistTrack, error) {
	var playlistTracks []spotifyAPI.PlaylistTrack
	cacheKey := "tracksPlaylist" + idPlaylist.String()

	ok, err := getJSON(cacheKey, &playlistTracks)
	if err != nil {
		return nil, err
	}

	if ok {
		log.Printf("Cache hit for playlist: %s", idPlaylist.String())
		return playlistTracks, nil
	}

	log.Printf("Cache miss for playlist: %s, fetching from API", idPlaylist.String())
	tracks, errGetAllPlaylistTracks := playlistService.GetAllPlaylistTracks(idPlaylist)
	if errGetAllPlaylistTracks != nil {
		return nil, errGetAllPlaylistTracks
	}

	_ = setJSON(cacheKey, tracks, defaultTTL)
	return tracks, nil
}

// GetCachedUserTracksOrSet retrieves user tracks from cache or API
func GetCachedUserTracksOrSet() ([]spotifyAPI.SavedTrack, error) {
	var userTracks []spotifyAPI.SavedTrack

	ok, err := getJSON("userTracks", &userTracks)
	if err != nil {
		return nil, err
	}
	if ok {
		log.Println("Cache hit for user tracks")
		return userTracks, nil
	}

	log.Println("Cache miss for user tracks, fetching from API")
	tracks, errTracks := userService.GetAllUserTracks()
	if errTracks != nil {
		return nil, errTracks
	}

	_ = setJSON("userTracks", tracks, defaultTTL)
	return tracks, nil
}
