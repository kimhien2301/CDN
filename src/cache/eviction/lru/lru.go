package lru

import "cache/eviction/modifiedlru"

func New(capacity int) *modifiedlru.CacheStorage {
	return modifiedlru.New(capacity, capacity)
}
