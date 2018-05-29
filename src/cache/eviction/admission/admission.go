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
	// cacheNum      int
}

func New(admissionList []interface{}) *CacheStorage {
	cache := new(CacheStorage)
	cache.Init()
	for _, key := range admissionList {
		cache.admissionList[key] = true
	}
	// fmt.Printf("Admission List: %v\n", cache.admissionList)
	return cache
}

func (cache *CacheStorage) CacheList() []interface{} {
	cacheList := make([]interface{}, 0)
	for key, _ := range cache.storage {
		cacheList = append(cacheList, key)
	}
	// fmt.Printf("Cache List: %v\n", cacheList)
	// fmt.Printf("Cache storage: %v\n", cache.storage)
	return cacheList
}

func (cache *CacheStorage) SetAdmissionList(admissionList []interface{}) {
	cache.admissionList = make(map[interface{}]bool)
	for _, key := range admissionList {
		cache.admissionList[key] = true
	}
	// fmt.Printf("Set Admission list of %d: %v\n", len(cache.admissionList), cache.admissionList)
}

func (cache *CacheStorage) Init() {
	cache.admissionList = make(map[interface{}]bool)
	cache.storage = make(map[interface{}]interface{})
	cache.ResetCount()
	// cache.cacheNum++
	// fmt.Println("Admission List")
	// fmt.Println(cache.admissionList)
	// fmt.Println("Storage")
	// fmt.Println(cache.storage)
}

func (cache *CacheStorage) Len() int {
	// fmt.Printf("Cache storage: %v\n", cache.storage)
	return len(cache.storage)
}

func (cache *CacheStorage) Exist(key interface{}) bool {
	_, exist := cache.storage[key]
	fmt.Printf("Cache storage with key[%v]: %v\n", cache.storage[key], exist)
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
		deleteCandidate := rand.Intn(len(deleteCandidates))
		fmt.Printf("Delete candidate %d of %d\n", deleteCandidate, len(deleteCandidates))
		delete(cache.storage, deleteCandidates[deleteCandidate])
		// delete(cache.storage, deleteCandidates[rand.Intn(len(deleteCandidates))])
	}
}

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	// fmt.Printf("Admission list with key[%d]: %v\n", key, cache.admissionList[key])
	if cache.admissionList[key] {
		// fmt.Printf("Storage with key[%d] before: %v\n", key, cache.storage[key])
		cache.storage[key] = value
		// fmt.Printf("Storage with key[%d] after: %v\n", key, cache.storage[key])
		cache.Evict()
		return value
	}
	return nil
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	content, exist := cache.storage[key]
	// fmt.Printf("Content %v: %v\n", content, exist)
	if exist {
		cache.hitCount++
		// fmt.Println("HIT")
		// cache.cacheNum++
		// if cache.cacheNum == 3 {
		// 	cache.cacheNum = 0
		// }
		return content
	}
	cache.missCount++
	// fmt.Println("MISS")
	// cache.cacheNum++
	// if cache.cacheNum == 3 {
	// 	cache.cacheNum = 0
	// }
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
