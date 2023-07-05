package cache

import (
	"time"

	"xarr-proxy/internal/config"

	cache "github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	gocache "github.com/patrickmn/go-cache"
)

var (
	c  *gocache.Cache
	cm *cache.Cache[string]
)

func Init(cfg *config.Config) {
	c = gocache.New(time.Duration(cfg.CacheTTL)*time.Second, time.Duration(cfg.CachePurgeInterval)*time.Second)
	gocacheStore := gocache_store.NewGoCache(c)
	cm = cache.New[string](gocacheStore)
}

func Get() *gocache.Cache {
	return c
}

func GetManager() *cache.Cache[string] {
	return cm
}
