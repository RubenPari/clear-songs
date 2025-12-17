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

var rdb *redis.Client
var ctx = context.Background()

const (
	defaultTTL = 5 * time.Minute
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

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	log.Println("Connected to Redis for caching")
}

// internal helpers
func setJSON(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, data, ttl).Err()
}

func getJSON(key string, target interface{}) (bool, error) {
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
	return setJSON("spotify_token", token, defaultTTL)
}

func GetToken() *oauth2.Token {
	var token oauth2.Token
	ok, err := getJSON("spotify_token", &token)
	if err != nil || !ok {
		return nil
	}
	return &token
}

// InvalidateUserData invalidates all user-related cache entries
func InvalidateUserData() {
	log.Println("Invalidating user data cache")
	_ = rdb.Del(ctx, "userTracks").Err()
	InvalidateAllPlaylists()
}

// InvalidatePlaylist invalidates cache for a specific playlist
func InvalidatePlaylist(playlistID spotifyAPI.ID) {
	log.Printf("Invalidating cache for playlist: %s", playlistID.String())
	_ = rdb.Del(ctx, "tracksPlaylist"+playlistID.String()).Err()
}

// InvalidateAllPlaylists invalidates all playlist caches
func InvalidateAllPlaylists() {
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
