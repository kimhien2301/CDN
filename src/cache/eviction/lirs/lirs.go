package lirs

import "container/list"
import "time"

type CacheStorage struct {
	capacity  int
	hitCount  int
	missCount int
	storage   *list.List
}

type Entry struct {
	key         interface{}
	value       interface{}
	accessTimes *list.List
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
	for element := cache.storage.Front(); element != nil; element = element.Next() {
		if element.Value.(Entry).key == key {
			return true
		}
	}
	return false
}

func (cache *CacheStorage) Evict() {
	var victimElement *list.Element
	var maxIRR time.Duration = 0
	for element := cache.storage.Back(); element != nil; element = element.Prev() {
		entry := element.Value.(Entry)
		if entry.accessTimes.Len() < 2 {
			victimElement = element
			break
		}
		irr := entry.accessTimes.Front().Value.(time.Time).Sub(entry.accessTimes.Back().Value.(time.Time))
		if maxIRR < irr {
			maxIRR = irr
			victimElement = element
		}
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
	entry := Entry{key, value, list.New()}
	entry.accessTimes.PushFront(time.Now())
	cache.storage.PushFront(entry)
	return value
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	if cache.Exist(key) {
		cache.hitCount++
		for element := cache.storage.Front(); element != nil; element = element.Next() {
			entry := element.Value.(Entry)
			if entry.key == key {
				entry.accessTimes.PushFront(time.Now())
				if entry.accessTimes.Len() > 2 {
					entry.accessTimes.Remove(entry.accessTimes.Back())
				}
				cache.storage.Remove(element)
				cache.storage.PushFront(entry)
				return entry.value
			}
		}
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
	cache.storage = list.New()
	cache.ResetCount()
}
