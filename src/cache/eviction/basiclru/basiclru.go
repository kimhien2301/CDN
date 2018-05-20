package basiclru

import "fmt"
import "container/list"

import "utils"

type Entry struct {
	key, value interface{}
}

type CacheStorage struct {
	capacity  int
	jump      int
	storage   *list.List
	hit, miss int
}

func New(capacity, jump int) *CacheStorage {
	cache := new(CacheStorage)
	cache.capacity = capacity
	cache.jump = jump
	cache.storage = list.New()
	cache.hit = 0
	cache.miss = 0
	return cache
}

func (cache *CacheStorage) CacheList() []interface{} {
	cacheList := make([]interface{}, 0)
	for element := cache.storage.Front(); element != nil; element = element.Next() {
		cacheList = append(cacheList, element.Value.(Entry).key)
	}
	return cacheList
}

func (cache *CacheStorage) Clear() {
	cache.storage = list.New()
	cache.ResetCount()
}

func (cache *CacheStorage) Capacity() int {
	return cache.capacity
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

func (cache *CacheStorage) Len() int {
	return cache.storage.Len()
}

func (cache *CacheStorage) Inspect() {
	for e := cache.storage.Front(); e != nil; e = e.Next() {
		utils.DebugPrint(fmt.Sprintln(e.Value))
	}
}

func (cache *CacheStorage) Init(capacity, jump int) *CacheStorage {
	if cache.storage == nil {
		cache.storage = list.New()
	}
	cache.capacity = capacity
	cache.jump = jump
	cache.storage.Init()
	return cache
}

func (cache *CacheStorage) evict() {
	for cache.Len() >= cache.capacity {
		cache.storage.Remove(cache.storage.Back())
	}
}

func (cache *CacheStorage) updatePosition(element *list.Element) {
	mark := element

	for index := 1; index < cache.jump; index++ {
		if mark == nil {
			break
		}
		mark = mark.Prev()
	}

	if mark != nil {
		cache.storage.MoveAfter(element, mark)
	} else {
		cache.storage.MoveToFront(element)
	}
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	fmt.Println("Content request: ", key)
	if cache.Exist(key) {
		fmt.Println("HIT")
		fmt.Printf("Cache storage[%v] before:\n", cache.storage.Len())
		for element := cache.storage.Front(); element != nil; element = element.Next() {
			fmt.Printf("KEY[%v]: %v\n", element.Value.(Entry).key, element.Value.(Entry).value)
		}

		cache.hit++
		for element := cache.storage.Front(); element != nil; element = element.Next() {
			if element.Value.(Entry).key == key {
				cache.updatePosition(element)
				return element.Value.(Entry).value
			}
		}

		fmt.Printf("Cache storage[%v] after:\n", cache.storage.Len())
		for element := cache.storage.Front(); element != nil; element = element.Next() {
			fmt.Printf("KEY[%v]: %v\n", element.Value.(Entry).key, element.Value.(Entry).value)
		}
	}
	fmt.Println("MISS")
	cache.miss++
	return nil
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
	if cache.Exist(key) || cache.capacity == 0 {
		return nil
	}

	fmt.Printf("Cache storage[%v] before insert:\n", cache.storage.Len())
	for element := cache.storage.Front(); element != nil; element = element.Next() {
		fmt.Printf("KEY[%v]: %v\n", element.Value.(Entry).key, element.Value.(Entry).value)
	}

	cache.evict()

	var element *list.Element
	entry := Entry{key, value}

	if cache.jump == 0 {
		element = cache.storage.PushBack(entry)
	} else {
		mark := cache.storage.Back()
		if mark != nil {
			for index := 1; index < cache.jump; index++ {
				if mark == nil {
					break
				}
				mark = mark.Prev()
			}
		}

		if mark == nil {
			element = cache.storage.PushFront(entry)
		} else {
			element = cache.storage.InsertAfter(entry, mark)
		}
	}

	fmt.Printf("Cache storage[%v] after insert:\n", cache.storage.Len())
	for element := cache.storage.Front(); element != nil; element = element.Next() {
		fmt.Printf("KEY[%v]: %v\n", element.Value.(Entry).key, element.Value.(Entry).value)
	}

	return element.Value.(Entry).value
}
