package arc

import "container/list"
import "math"

type CacheStorage struct {
	t1        *list.List
	b1        *list.List
	t2        *list.List
	b2        *list.List
	p         int
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
	cache.t1 = list.New()
	cache.b1 = list.New()
	cache.t2 = list.New()
	cache.b2 = list.New()
	cache.capacity = capacity
	cache.p = cache.capacity / 2
	return cache
}

func (cache *CacheStorage) CacheList() []interface{} {
	cacheList := make([]interface{}, 0)
	for element := cache.t1.Front(); element != nil; element = element.Next() {
		cacheList = append(cacheList, element.Value.(Entry).key)
	}
	for element := cache.t2.Front(); element != nil; element = element.Next() {
		cacheList = append(cacheList, element.Value.(Entry).key)
	}
	return cacheList
}

func (cache *CacheStorage) Len() int {
	return cache.t1.Len() + cache.t2.Len()
}

func (cache *CacheStorage) Capacity() int {
	return cache.capacity
}

func (cache *CacheStorage) ExistInT1(key interface{}) bool {
	for element := cache.t1.Front(); element != nil; element = element.Next() {
		if element.Value.(Entry).key == key {
			return true
		}
	}
	return false
}

func (cache *CacheStorage) ExistInT2(key interface{}) bool {
	for element := cache.t2.Front(); element != nil; element = element.Next() {
		if element.Value.(Entry).key == key {
			return true
		}
	}
	return false
}

func (cache *CacheStorage) ExistInB1(key interface{}) bool {
	for element := cache.b1.Front(); element != nil; element = element.Next() {
		if element.Value.(Entry).key == key {
			return true
		}
	}
	return false
}

func (cache *CacheStorage) ExistInB2(key interface{}) bool {
	for element := cache.b2.Front(); element != nil; element = element.Next() {
		if element.Value.(Entry).key == key {
			return true
		}
	}
	return false
}

func (cache *CacheStorage) ExistInAny(key interface{}) bool {
	if cache.ExistInT1(key) {
		return true
	} else if cache.ExistInT2(key) {
		return true
	} else if cache.ExistInB1(key) {
		return true
	} else if cache.ExistInB2(key) {
		return true
	}
	return false
}

func (cache *CacheStorage) Exist(key interface{}) bool {
	if cache.ExistInT1(key) {
		return true
	} else if cache.ExistInT2(key) {
		return true
	}
	return false
}

func (cache *CacheStorage) Replace(key interface{}) {
	if cache.t1.Len() > 0 && (cache.t1.Len() > cache.p || (cache.ExistInB2(key) && cache.t1.Len() == cache.p)) {
		lruPage := cache.t1.Back()
		cache.t1.Remove(lruPage)
		cache.b1.PushFront(lruPage.Value.(Entry))
	} else {
		lruPage := cache.t2.Back()
		cache.t2.Remove(lruPage)
		cache.b2.PushFront(lruPage.Value.(Entry))
	}
}

func (cache *CacheStorage) PickElementInB1(key interface{}) *list.Element {
	for element := cache.b1.Front(); element != nil; element = element.Next() {
		if element.Value.(Entry).key == key {
			return element
		}
	}
	return nil
}

func (cache *CacheStorage) PickElementInB2(key interface{}) *list.Element {
	for element := cache.b2.Front(); element != nil; element = element.Next() {
		if element.Value.(Entry).key == key {
			return element
		}
	}
	return nil
}

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	if cache.ExistInB1(key) {
		// cache.missCount++
		var delta int
		if cache.b1.Len() >= cache.b2.Len() {
			delta = 1
		} else {
			delta = cache.b2.Len() / cache.b1.Len()
		}
		cache.p = int(math.Min(float64(cache.p+delta), float64(cache.capacity)))
		cache.Replace(key)
		cache.b1.Remove(cache.PickElementInB1(key))
		cache.t2.PushFront(Entry{key, value})
		return value
	} else if cache.ExistInB2(key) {
		var delta int
		if cache.b2.Len() >= cache.b1.Len() {
			delta = 1
		} else {
			delta = cache.b1.Len() / cache.b2.Len()
		}
		cache.p = int(math.Max(float64(cache.p-delta), 0))
		cache.Replace(key)
		cache.b2.Remove(cache.PickElementInB2(key))
		cache.t2.PushFront(Entry{key, value})
	} else if !cache.ExistInAny(key) {
		if cache.t1.Len()+cache.b1.Len() == cache.capacity {
			if cache.t1.Len() < cache.capacity {
				cache.b1.Remove(cache.b1.Back())
				cache.Replace(key)
			} else {
				cache.t1.Remove(cache.t1.Back())
			}
		} else if cache.t1.Len()+cache.b1.Len() < cache.capacity {
			if cache.t1.Len()+cache.t2.Len()+cache.b1.Len()+cache.b2.Len() >= cache.capacity {
				if cache.t1.Len()+cache.t2.Len()+cache.b1.Len()+cache.b2.Len() == 2*cache.capacity {
					cache.b2.Remove(cache.b2.Back())
				}
				cache.Replace(key)
			}
		}
		cache.t1.PushFront(Entry{key, value})
	}
	return value
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	var element *list.Element
	var entry Entry

	if cache.ExistInT1(key) {
		cache.hitCount++
		for element = cache.t1.Front(); element != nil; element = element.Next() {
			if element.Value.(Entry).key == key {
				break
			}
		}
		entry = element.Value.(Entry)
		cache.t1.Remove(element)
		cache.t2.PushFront(entry)
		return entry.value
	} else if cache.ExistInT2(key) {
		cache.hitCount++
		for element = cache.t2.Front(); element != nil; element = element.Next() {
			if element.Value.(Entry).key == key {
				break
			}
		}
		entry = element.Value.(Entry)
		cache.t2.MoveToFront(element)
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
	cache.t1 = list.New()
	cache.t2 = list.New()
	cache.b1 = list.New()
	cache.b2 = list.New()
}
