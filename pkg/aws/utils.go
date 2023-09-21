package aws

import "sync"

type CacheKey string

var (
	snsCacheKey     CacheKey = "sns"
	sesCacheKey     CacheKey = "ses"
	sessionCacheKey CacheKey = "session"
	sqsCacheKey     CacheKey = "sqs"
)

var cache = make(map[CacheKey]map[string]any)
var mu sync.Mutex

func addServiceToCache(key CacheKey, region string, svc any) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := cache[key]; !exists {
		cache[key] = make(map[string]any)
	}
	cache[key][region] = svc
}

func getServiceFromCache(key CacheKey, region string) any {
	if cached, exists := cache[key][region]; exists {
		return cached
	}
	return nil
}

func String(s string) *string {
	return &s
}
