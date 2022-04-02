package user

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	defaultExpiration = 5 * time.Minute
	cleanupInterval   = 10 * time.Minute
)

var userCache *cache.Cache

func init() {
	userCache = cache.New(defaultExpiration, cleanupInterval)
}

func addToCache(id string, user User) {
	userCache.Set(id, user, cache.DefaultExpiration)
}
func getFromCache(id string) *User {
	if x, found := userCache.Get(id); found {
		user := x.(User)
		return &user
	}
	return nil
}
