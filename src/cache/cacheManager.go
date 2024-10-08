package cache

import (
	"time"

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
