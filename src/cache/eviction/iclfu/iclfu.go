package iclfu

import "math/rand"

type CacheStorage struct {
	storage   map[interface{}]Entry
	capacity  int
	hitCount  int
	missCount int
}

type Entry struct {
	key     interface{}
	value   interface{}
	counter uint64
}

func New(capacity int) *CacheStorage {
	cache := new(CacheStorage)
	cache.storage = make(map[interface{}]Entry)
	cache.capacity = capacity
	cache.ResetCount()
	return cache
}

func (cache *CacheStorage) CacheList() []interface{} {
	cacheList := make([]interface{}, 0)
	for key, _ := range cache.storage {
		cacheList = append(cacheList, key)
	}
	return cacheList
}

func (cache *CacheStorage) Clear() {
	cache.storage = make(map[interface{}]Entry)
	cache.ResetCount()
}

func (cache *CacheStorage) Len() int {
	return len(cache.storage)
}

func (cache *CacheStorage) Exist(key interface{}) bool {
	_, exist := cache.storage[key]
	return exist
}

func (cache *CacheStorage) evictLFUEntry() {
	var minCounter uint64
	minCounter--
	for _, entry := range cache.storage {
		if minCounter > entry.counter {
			minCounter = entry.counter
		}
	}
	deleteCandidate := make([]interface{}, 0)
	for key, entry := range cache.storage {
		if entry.counter == minCounter {
			deleteCandidate = append(deleteCandidate, key)
		}
	}

	delete(cache.storage, deleteCandidate[rand.Intn(len(deleteCandidate))])
}

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	if !cache.Exist(key) {
		for cache.capacity <= len(cache.storage) {
			cache.evictLFUEntry()
		}
		cache.storage[key] = Entry{key, value, 0}
		return value
	}
	return nil
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	entry, exist := cache.storage[key]
	if exist {
		cache.hitCount++
		entry.counter++
		return entry.value
	}
	cache.missCount++
	return nil
}

func (cache *CacheStorage) Capacity() int {
	return cache.capacity
}

func (cache *CacheStorage) HitCount() int {
	return cache.hitCount
}

func (cache *CacheStorage) MissCount() int {
	return cache.missCount
}

func (cache *CacheStorage) ResetCount() {
	cache.hitCount = 0
	cache.missCount = 0
}
