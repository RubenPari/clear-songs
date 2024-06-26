package cacheManager

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Cache *cache.Cache

// Init creates a cacheManager with a default
// expiration time of 5 and which purges
// expired items every 10
func Init() {
	Cache = cache.New(5*time.Minute, 10*time.Minute)
}

// Delete flushes the cacheManager
func Delete() {
	Cache.Flush()
}

// Set adds a new item to the cache
// The item is added with the default expiration time
func Set(key string, value interface{}) {
	Cache.Set(key, value, cache.DefaultExpiration)
}

// Get retrieves an item from the cache
// If the item does not exist, it returns nil
func Get(key string) interface{} {
	value, found := Cache.Get(key)
	if found {
		return value
	}
	return nil
}
