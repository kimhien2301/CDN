package admission

import (
	"fmt"
	"math/rand"
)

type CacheStorage struct {
	admissionList map[interface{}]bool
	storage       map[interface{}]interface{}
	hitCount      int
	missCount     int
}

func New(admissionList []interface{}) *CacheStorage {
	cache := new(CacheStorage)
	cache.Init()

	fmt.Println("Admission List")
	fmt.Println(cache.admissionList)

	for _, key := range admissionList {
		cache.admissionList[key] = true
	}
	fmt.Println("Admission List")
	fmt.Println(cache.admissionList)
	return cache
}

func (cache *CacheStorage) CacheList() []interface{} {
	cacheList := make([]interface{}, 0)
	for key, _ := range cache.storage {
		cacheList = append(cacheList, key)
	}
	// fmt.Println("Cache List")
	// fmt.Println(cacheList)
	return cacheList
}

func (cache *CacheStorage) SetAdmissionList(admissionList []interface{}) {
	cache.admissionList = make(map[interface{}]bool)
	for _, key := range admissionList {
		cache.admissionList[key] = true
	}
	// fmt.Println(admissionList)
}

func (cache *CacheStorage) Init() {
	cache.admissionList = make(map[interface{}]bool)
	cache.storage = make(map[interface{}]interface{})
	cache.ResetCount()
	// fmt.Println("Admission List")
	// fmt.Println(cache.admissionList)
	// fmt.Println("Storage")
	// fmt.Println(cache.storage)
}

func (cache *CacheStorage) Len() int {
	return len(cache.storage)
}

func (cache *CacheStorage) Exist(key interface{}) bool {
	_, exist := cache.storage[key]
	return exist
}

func (cache *CacheStorage) Evict() {
	for len(cache.storage) > len(cache.admissionList) {
		deleteCandidates := make([]interface{}, 0)
		for key := range cache.storage {
			if cache.admissionList[key] {
				deleteCandidates = append(deleteCandidates, key)
			}
		}
		delete(cache.storage, deleteCandidates[rand.Intn(len(deleteCandidates))])
	}
}

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	if cache.admissionList[key] {
		cache.storage[key] = value
		cache.Evict()
		return value
	}
	return nil
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	content, exist := cache.storage[key]
	if exist {
		cache.hitCount++
		return content
	}
	cache.missCount++
	return nil
}

func (cache *CacheStorage) Capacity() int {
	return len(cache.admissionList)
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
	cache.storage = make(map[interface{}]interface{})
	cache.ResetCount()
}
