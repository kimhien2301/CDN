package modifiedlru

import "fmt"
import "container/list"
import "sync"
import "utils"

type Entry struct {
	key, value interface{}
}

type CacheStorage struct {
	capacity  int
	jump      int
	storage   *list.List
	mutex     sync.Mutex
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
	for e := cache.front(); e != nil; e = e.Next() {
		utils.DebugPrint(fmt.Sprintln(e.Value))
	}
}

func (cache *CacheStorage) back() *list.Element {
	return cache.storage.Back()
}

func (cache *CacheStorage) front() *list.Element {
	return cache.storage.Front()
}

func (cache *CacheStorage) remove(element *list.Element) interface{} {
	cache.mutex.Lock()
	value := cache.storage.Remove(element)
	cache.mutex.Unlock()
	return value
}

func (cache *CacheStorage) insertAfter(value interface{}, mark *list.Element) *list.Element {
	cache.mutex.Lock()
	element := cache.storage.InsertAfter(value, mark)
	cache.mutex.Unlock()
	return element
}

func (cache *CacheStorage) insertBefore(value interface{}, mark *list.Element) *list.Element {
	cache.mutex.Lock()
	element := cache.storage.InsertBefore(value, mark)
	cache.mutex.Unlock()
	return element
}

func (cache *CacheStorage) moveAfter(element, mark *list.Element) {
	cache.mutex.Lock()
	cache.storage.MoveAfter(element, mark)
	cache.mutex.Unlock()
}

func (cache *CacheStorage) moveBefore(element, mark *list.Element) {
	cache.mutex.Lock()
	cache.storage.MoveBefore(element, mark)
	cache.mutex.Unlock()
}

func (cache *CacheStorage) moveToFront(element *list.Element) {
	cache.mutex.Lock()
	cache.storage.MoveToFront(element)
	cache.mutex.Unlock()
}

func (cache *CacheStorage) moveToBack(element *list.Element) {
	cache.mutex.Lock()
	cache.storage.MoveToBack(element)
	cache.mutex.Unlock()
}

func (cache *CacheStorage) pushBack(value interface{}) *list.Element {
	cache.mutex.Lock()
	element := cache.storage.PushBack(value)
	cache.mutex.Unlock()
	return element
}

func (cache *CacheStorage) pushFront(value interface{}) *list.Element {
	cache.mutex.Lock()
	element := cache.storage.PushFront(value)
	cache.mutex.Unlock()
	return element
}

func (cache *CacheStorage) Init(capacity, jump int) *CacheStorage {
	cache.mutex.Lock()
	if cache.storage == nil {
		cache.storage = list.New()
	}
	cache.capacity = capacity
	cache.jump = jump
	cache.storage.Init()
	cache.mutex.Unlock()
	return cache
}

func (cache *CacheStorage) evict() {
	for cache.Len() >= cache.capacity {
		cache.remove(cache.back())
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
		cache.moveAfter(element, mark)
	} else {
		cache.moveToFront(element)
	}
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	if cache.Exist(key) {
		cache.hit++
		for element := cache.front(); element != nil; element = element.Next() {
			if element.Value.(Entry).key == key {
				cache.updatePosition(element)
				return element.Value.(Entry).value
			}
		}
	}
	cache.miss++
	return nil
}

func (cache *CacheStorage) Exist(key interface{}) bool {
	for element := cache.front(); element != nil; element = element.Next() {
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

	cache.evict()

	var element *list.Element
	entry := Entry{key, value}

	if cache.jump == 0 {
		element = cache.pushBack(entry)
	} else {
		mark := cache.back()
		if mark != nil {
			for index := 1; index < cache.jump; index++ {
				if mark == nil {
					break
				}
				mark = mark.Prev()
			}
		}

		if mark == nil {
			element = cache.pushFront(entry)
		} else {
			element = cache.insertAfter(entry, mark)
		}
	}
	return element.Value.(Entry).value
}
