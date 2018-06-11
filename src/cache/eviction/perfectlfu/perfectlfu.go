// TU add
package perfectlfu

import (
	"container/list"
)

type Entry struct {
	key, value interface{}
}

type CacheStorage struct {
	capacity  int
	storage   *list.List
	hit, miss int
	frequency map[interface{}]int
}

func New(capacity int) *CacheStorage {
	cache := new(CacheStorage)
	cache.capacity = capacity
	cache.storage = list.New()
	cache.frequency = make(map[interface{}]int)
	cache.hit = 0
	cache.miss = 0
	return cache
}

func (cache *CacheStorage) Len() int {
	return cache.storage.Len()
}

func (cache *CacheStorage) Capacity() int {
	return cache.capacity
}

func (cache *CacheStorage) CacheList() []interface{} {
	cacheList := make([]interface{}, 0)
	for element := cache.storage.Front(); element != nil; element = element.Next() {
		cacheList = append(cacheList, element.Value.(Entry).key)
	}
	return cacheList
}

func (cache *CacheStorage) Exist(key interface{}) bool {
	for element := cache.storage.Front(); element != nil; element = element.Next() {
		if element.Value.(Entry).key == key {
			return true
		}
	}
	return false
}

//
// func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
// 	if cache.Exist(key) || cache.capacity == 0 {
// 		return nil
// 	}

// 	cache.Evict()
// 	// First time inserted to cache
// 	cache.storage.PushBack(Entry{key, value})
// 	cache.frequency[key] = 1
// 	return value
// }

// //
// func (cache *CacheStorage) Evict() {
// 	for cache.Len() >= cache.capacity {
// 		cache.remove()
// 	}

// }

func (cache *CacheStorage) EvictWithKey(key interface{}) {

	for element := cache.storage.Front(); element != nil; element = element.Next() {
		if element.Value.(Entry).key == key {
			delete(cache.frequency, key)
			cache.storage.Remove(element)
			return
		}
	}
}

//
// func (cache *CacheStorage) remove() {
// 	removedEle := cache.storage.Front()
// 	for element := cache.storage.Front().Next(); element != nil; element = element.Next() {
// 		if cache.frequency[element.Value.(Entry).key] < cache.frequency[removedEle.Value.(Entry).key] {
// 			removedEle = element
// 		}
// 	}
// 	delete(cache.frequency, removedEle.Value.(Entry).key)
// 	cache.storage.Remove(removedEle)
// }

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	if cache.Exist(key) {

		cache.hit++
		for element := cache.storage.Front(); element != nil; element = element.Next() {
			if element.Value.(Entry).key == key {
				cache.frequency[key]++
				return element.Value.(Entry).value
			}
		}
		return 0
	}
	cache.miss++
	return nil
}

func (cache *CacheStorage) HitCount() int {
	return cache.hit
}

func (cache *CacheStorage) MissCount() int {
	return cache.miss
}

func (cache *CacheStorage) ResetCount() {
	cache.hit = 0
	cache.miss = 0
}

func (cache *CacheStorage) Clear() {
	cache.storage = list.New()
	cache.frequency = make(map[interface{}]int)
	cache.ResetCount()
}

/////////////////////////

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	if cache.Exist(key) || cache.capacity == 0 {
		return nil
	}

	_, exist := cache.frequency[key]

	if exist {
		cache.frequency[key]++
	} else {
		cache.frequency[key] = 1
	}

	//cache.Evict()

	if cache.Len() < cache.capacity {
		cache.storage.PushBack(Entry{key, value})
	} else {
		// get content with minimum # of access
		removedEle := cache.storage.Front()
		for element := cache.storage.Front().Next(); element != nil; element = element.Next() {
			if cache.frequency[element.Value.(Entry).key] < cache.frequency[removedEle.Value.(Entry).key] {
				removedEle = element
			}
		}

		// Decide if server should cache new content
		if cache.frequency[key] > cache.frequency[removedEle.Value.(Entry).key] {
			cache.storage.Remove(removedEle)
			cache.storage.PushBack(Entry{key, value})
		} else {
			return nil
		}
	}

	return value
}
