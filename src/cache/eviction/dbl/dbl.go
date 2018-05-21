package dbl

import "container/list"

type CacheStorage struct {
	l1        *list.List
	l2        *list.List
	capacity  int
	hitCount  int
	missCount int
}

type Entry struct {
	key   interface{}
	value interface{}
}

func New(capacity int) *CacheStorage {
	cache := new(CacheStorage)
	cache.l1 = list.New()
	cache.l2 = list.New()
	cache.capacity = capacity
	return cache
}

func (cache *CacheStorage) CacheList() []interface{} {
	cacheList := make([]interface{}, 0)
	for element := cache.l1.Front(); element != nil; element = element.Next() {
		cacheList = append(cacheList, element.Value.(Entry).key)
	}
	for element := cache.l2.Front(); element != nil; element = element.Next() {
		cacheList = append(cacheList, element.Value.(Entry).key)
	}
	return cacheList
}

func (cache *CacheStorage) Len() int {
	return cache.l1.Len() + cache.l2.Len()
}

func (cache *CacheStorage) Capacity() int {
	return cache.capacity
}

func (cache *CacheStorage) ExistInL1(key interface{}) bool {
	for element := cache.l1.Front(); element != nil; element = element.Next() {
		if element.Value.(Entry).key == key {
			return true
		}
	}
	return false
}

func (cache *CacheStorage) ExistInL2(key interface{}) bool {
	for element := cache.l2.Front(); element != nil; element = element.Next() {
		if element.Value.(Entry).key == key {
			return true
		}
	}
	return false
}

func (cache *CacheStorage) ExistInAny(key interface{}) bool {
	if cache.ExistInL1(key) {
		return true
	} else if cache.ExistInL2(key) {
		return true
	}
	return false
}

func (cache *CacheStorage) Exist(key interface{}) bool {
	if cache.ExistInL1(key) {
		return true
	} else if cache.ExistInL2(key) {
		return true
	}
	return false
}

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	if cache.l1.Len() == cache.capacity/2 {
		cache.l1.Remove(cache.l1.Back())
	} else if cache.l1.Len() < cache.capacity/2 {
		if cache.l1.Len()+cache.l2.Len() == cache.capacity {
			cache.l2.Remove(cache.l2.Back())
		}
	}
	cache.l1.PushFront(Entry{key, value})
	return value
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	var element *list.Element
	var entry Entry

	if cache.ExistInL1(key) {
		cache.hitCount++
		for element = cache.l1.Front(); element != nil; element = element.Next() {
			if element.Value.(Entry).key == key {
				break
			}
		}
		entry = element.Value.(Entry)
		cache.l1.Remove(element)
		cache.l2.PushFront(entry)
		return entry.value
	} else if cache.ExistInL2(key) {
		cache.hitCount++
		for element = cache.l2.Front(); element != nil; element = element.Next() {
			if element.Value.(Entry).key == key {
				break
			}
		}
		entry = element.Value.(Entry)
		cache.l2.MoveToFront(element)
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
	cache.l1 = list.New()
	cache.l2 = list.New()
}
