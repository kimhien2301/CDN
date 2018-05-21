package windowlfu

import (
	"algorithm"
	"cache/eviction/admission"
	"fmt"
	"sort"
)

type CacheStorage struct {
	admissionCache *admission.CacheStorage
	capacity       int
	window         int
	history        []interface{}
}

func New(capacity, window int) *CacheStorage {
	wlfu := new(CacheStorage)
	admissionList := make([]interface{}, 0)
	wlfu.admissionCache = admission.New(admissionList)
	wlfu.capacity = capacity
	wlfu.window = window
	wlfu.history = make([]interface{}, 0)
	// fmt.Printf("Admission cache of %d: %v\nCapacity: %v\nWindow: %v\nHistory: %v\n", len(admissionList), wlfu.admissionCache, wlfu.capacity, wlfu.window, wlfu.history)
	return wlfu
}

func (cache *CacheStorage) CacheList() []interface{} {
	// fmt.Printf("Cache list: %v\n", cache.admissionCache.CacheList())
	return cache.admissionCache.CacheList()
}

func (cache *CacheStorage) Clear() {
	// fmt.Println("Admission Cache Clear")
	cache.admissionCache.Clear()
}

func (cache *CacheStorage) Len() int {
	// fmt.Printf("Admission Cache Lenght: %d\n", cache.admissionCache.Len())
	return cache.admissionCache.Len()
}

func (cache *CacheStorage) Exist(key interface{}) bool {
	// fmt.Printf("Key %v: %v\n", key, cache.admissionCache.Exist(key))
	return cache.admissionCache.Exist(key)
}

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	fmt.Printf("Insert key %v with value %v\n", key, value)
	return cache.admissionCache.Insert(key, value)
}

func (cache *CacheStorage) updateAdmissionList() {
	popularity := make(map[interface{}]int)
	for _, key := range cache.history {
		value, exist := popularity[key]
		if exist {
			popularity[key] = value + 1
		} else {
			popularity[key] = 1
		}
	}

	list := algorithm.List{}
	for key, counter := range popularity {
		entry := algorithm.Entry{key, float64(counter)}
		list = append(list, entry)
	}
	sort.Sort(list)
	fmt.Printf("List after sorted [%d]: %v\n", len(list), list)

	admissionList := make([]interface{}, 0)

	for i := len(list) - 1; i >= 0; i-- {
		admissionList = append(admissionList, list[i].Key)
	}

	cache.admissionCache.SetAdmissionList(admissionList[0:100])
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	cache.history = append(cache.history, key)
	content := cache.admissionCache.Fetch(key)
	// fmt.Printf("Admission cache: %v\n", cache.admissionCache)
	if len(cache.history) == cache.window {
		cache.updateAdmissionList()
		fmt.Println("Update admission list")
		fmt.Printf("Cache history [%d]: %v\n", len(cache.history), cache.history)
		cache.history = make([]interface{}, 0)
	}
	return content
}

func (cache *CacheStorage) Capacity() int {
	return cache.admissionCache.Capacity()
}

func (cache *CacheStorage) HitCount() int {
	return cache.admissionCache.HitCount()
}

func (cache *CacheStorage) MissCount() int {
	return cache.admissionCache.MissCount()
}

func (cache *CacheStorage) ResetCount() {
	cache.admissionCache.ResetCount()
}
