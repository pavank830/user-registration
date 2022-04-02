package user

import (
	"fmt"
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
	fmt.Println("added to cache id", id)
	userCache.Set(id, user, cache.DefaultExpiration)
}
func getFromCache(id string) *User {
	fmt.Println("got from cache id 1", id)
	if x, found := userCache.Get(id); found {
		fmt.Println("got from cache id 2", id)
		user := x.(User)
		return &user
	}
	return nil
}
