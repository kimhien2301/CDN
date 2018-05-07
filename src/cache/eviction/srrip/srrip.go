package srrip

import "container/list"
import "fmt"
import "utils"

type CacheStorage struct {
	storage   *list.List
	capacity  int
	rrpvMax   int
	hitCount  int
	missCount int
}

type Entry struct {
	key   interface{}
	value interface{}
	rrpv  int
}

func New(capacity, rrpvbit int) *CacheStorage {
	cache := new(CacheStorage)
	cache.capacity = capacity
	cache.storage = list.New()
	cache.rrpvMax = 1
	for i := 0; i < rrpvbit; i++ {
		cache.rrpvMax <<= 1
	}
	cache.rrpvMax--
	for index := 0; index < cache.capacity; index++ {
		cache.storage.PushFront(Entry{nil, nil, cache.rrpvMax})
	}
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

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	for {
		for element := cache.storage.Front(); element != nil; element = element.Next() {
			if element.Value.(Entry).rrpv == cache.rrpvMax {
				element.Value = Entry{key, value, cache.rrpvMax - 1}
				return value
			}
		}
		for element := cache.storage.Front(); element != nil; element = element.Next() {
			element.Value = Entry{element.Value.(Entry).key, element.Value.(Entry).value, element.Value.(Entry).rrpv + 1}
		}
	}
	return nil
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	if cache.Exist(key) {
		cache.hitCount++
		for element := cache.storage.Front(); element != nil; element = element.Next() {
			if element.Value.(Entry).key == key {
				entry := element.Value.(Entry)
				entry.rrpv = 0
				element.Value = entry
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
	cache.ResetCount()
	cache.storage = list.New()
}

func (cache *CacheStorage) Inspect() {
	utils.DebugPrint(fmt.Sprint("{ "))
	for element := cache.storage.Front(); element != nil; element = element.Next() {
		utils.DebugPrint(fmt.Sprintf("(%v, %d) ", element.Value.(Entry).key, element.Value.(Entry).rrpv))
	}
	utils.DebugPrint(fmt.Sprintln("}"))
}
