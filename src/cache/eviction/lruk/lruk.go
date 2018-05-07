package lruk

import "time"
import "container/list"

type CacheStorage struct {
	k         int
	capacity  int
	hitCount  int
	missCount int
	storage   map[interface{}]Entry
}

type Entry struct {
	value       interface{}
	accessTimes *list.List
}

func New(capacity, k int) *CacheStorage {
	cache := new(CacheStorage)
	cache.storage = make(map[interface{}]Entry)
	cache.capacity = capacity
	cache.k = k
	return cache
}

func (cache *CacheStorage) CacheList() []interface{} {
	cacheList := make([]interface{}, 0)
	for key, _ := range cache.storage {
		cacheList = append(cacheList, key)
	}
	return cacheList
}

func (cache *CacheStorage) Len() int {
	return len(cache.storage)
}

func (cache *CacheStorage) Capacity() int {
	return cache.capacity
}

func (cache *CacheStorage) Exist(key interface{}) bool {
	_, exist := cache.storage[key]
	return exist
}

func (cache *CacheStorage) Evict() {
	var oldestEntryKey interface{}

	for k := 1; k <= cache.k; k++ {
		oldestAccessTime := time.Now()
		for key, entry := range cache.storage {
			if entry.accessTimes.Len() == k {
				oldAccessTime := entry.accessTimes.Back()
				if oldAccessTime.Value.(time.Time).Before(oldestAccessTime) {
					oldestAccessTime = oldAccessTime.Value.(time.Time)
					oldestEntryKey = key
				}
			}
		}
		if oldestEntryKey != nil {
			break
		}
	}

	delete(cache.storage, oldestEntryKey)
}

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	if cache.Exist(key) {
		return nil
	}
	for cache.Len() >= cache.Capacity() {
		cache.Evict()
	}

	accessTime := time.Now()
	entry := Entry{value, list.New()}
	entry.accessTimes.PushFront(accessTime)

	cache.storage[key] = entry
	return value
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	if cache.Exist(key) {
		cache.hitCount++
		entry := cache.storage[key]
		accessTime := time.Now()
		entry.accessTimes.PushFront(accessTime)
		for entry.accessTimes.Len() > cache.k {
			entry.accessTimes.Remove(entry.accessTimes.Back())
		}
		return entry.value
	}
	cache.missCount++
	return nil
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

func (cache *CacheStorage) Clear() {
	cache.ResetCount()
	cache.storage = make(map[interface{}]Entry)
}
