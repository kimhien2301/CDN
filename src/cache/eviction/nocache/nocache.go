package nocache

import "math/rand"
import "container/list"

type CacheStorage struct {
	capacity  int
	hitCount  int
	missCount int
	storage   *list.List
}

type Entry struct {
	key   interface{}
	value interface{}
}

func New(capacity int) *CacheStorage {
	cache := new(CacheStorage)
	cache.capacity = capacity
	cache.storage = list.New()
	return cache
}

func (cache *CacheStorage) CacheList() []interface{} {
	cacheList := make([]interface{}, 0)
	for element := cache.storage.Front(); element != nil; element = element.Next() {
		cacheList = append(cacheList, element.Value.(Entry).key)
	}
	return cacheList
}

func (cache *CacheStorage) Len() int {
	return cache.storage.Len()
}

func (cache *CacheStorage) Capacity() int {
	return cache.capacity
}

func (cache *CacheStorage) Exist(key interface{}) bool {
	// for element := cache.storage.Front(); element != nil; element = element.Next() {
	// 	if element.Value.(Entry).key == key {
	// 		return true
	// 	}
	// }
	return false
}

func (cache *CacheStorage) Evict() {
	index := rand.Intn(cache.capacity)
	victimElement := cache.storage.Front()
	for i := 1; i <= index; i++ {
		victimElement = victimElement.Next()
	}
	cache.storage.Remove(victimElement)
}

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	if cache.Exist(key) {
		return nil
	}
	for cache.storage.Len() >= cache.capacity {
		cache.Evict()
	}
	cache.storage.PushFront(Entry{key, value})
	return value
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	// if cache.Exist(key) {
	// 	cache.hitCount++
	// 	for element := cache.storage.Front(); element != nil; element = element.Next() {
	// 		if element.Value.(Entry).key == key {
	// 			return element.Value.(Entry).value
	// 		}
	// 	}
	// }
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
	cache.storage = list.New()
}
