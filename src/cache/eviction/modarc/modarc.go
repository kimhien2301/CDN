package modarc

import (
	"container/list"
	"fmt"
	"math"
)

type CacheStorage struct {
	t1        *list.List
	b1        *list.List
	t2        *list.List
	b2        *list.List
	p         int
	capacity  int
	jump      int
	hitCount  int
	missCount int
	nodeID    string
}

type Entry struct {
	key   interface{}
	value interface{}
}

func New(capacity, jump int, node string) *CacheStorage {
	cache := new(CacheStorage)
	cache.t1 = list.New()
	cache.b1 = list.New()
	cache.t2 = list.New()
	cache.b2 = list.New()
	cache.capacity = capacity
	cache.jump = jump
	cache.p = cache.capacity / 2
	cache.nodeID = node
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

func (cache *CacheStorage) updatePosition(list *list.List, element *list.Element) {
	mark := element

	for index := 1; index < cache.jump; index++ {
		if mark == nil {
			break
		}
		mark = mark.Prev()
	}

	if mark != nil {
		list.MoveAfter(element, mark)
	} else {
		list.MoveToFront(element)
	}
}

func (cache *CacheStorage) pushElement(l *list.List, key, value interface{}) interface{} {
	var element *list.Element
	entry := Entry{key, value}
	if cache.jump == 0 {
		element = l.PushBack(entry)
	} else {
		mark := cache.t2.Back()
		if mark != nil {
			for index := 1; index < cache.jump; index++ {
				if mark == nil {
					break
				}
				mark = mark.Prev()
			}
		}

		if mark == nil {
			element = cache.t2.PushFront(entry)
		} else {
			element = cache.t2.InsertAfter(entry, mark)
		}
	}
}

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	fmt.Println("INSERT")
	fmt.Println("Cache server: ", cache.nodeID)
	if cache.ExistInB1(key) {
		cache.missCount++
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

		fmt.Printf("Cache T1[%v] before:\n", cache.t1.Len())
		for e := cache.t1.Front(); e != nil; e = e.Next() {
			fmt.Printf("KEY[%v]: %v\n", e.Value.(Entry).key, e.Value.(Entry).value)
		}

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

		fmt.Printf("Cache T1[%v] after:\n", cache.t1.Len())
		for e := cache.t1.Front(); e != nil; e = e.Next() {
			fmt.Printf("KEY[%v]: %v\n", e.Value.(Entry).key, e.Value.(Entry).value)
		}
	}
	fmt.Println("--------------------")
	return value
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	var element *list.Element
	var entry Entry
	fmt.Println("FETCH")
	fmt.Println("Cache server: ", cache.nodeID)
	fmt.Printf("Request content: %v\n", key)
	if cache.ExistInT1(key) {
		fmt.Printf("HIT in T1\nCache T1[%v] before:\n", cache.t1.Len())
		for e := cache.t1.Front(); e != nil; e = e.Next() {
			fmt.Printf("KEY[%v]: %v\n", e.Value.(Entry).key, e.Value.(Entry).value)
		}
		fmt.Printf("Cache T2[%v] before:\n", cache.t2.Len())
		for e := cache.t2.Front(); e != nil; e = e.Next() {
			fmt.Printf("KEY[%v]: %v\n", e.Value.(Entry).key, e.Value.(Entry).value)
		}

		cache.hitCount++
		for element = cache.t1.Front(); element != nil; element = element.Next() {
			if element.Value.(Entry).key == key {
				break
			}
		}

		entry = element.Value.(Entry)
		cache.t1.Remove(element)

		if cache.jump == 0 {
			element = cache.t2.PushBack(entry)
		} else {
			mark := cache.t2.Back()
			if mark != nil {
				for index := 1; index < cache.jump; index++ {
					if mark == nil {
						break
					}
					mark = mark.Prev()
				}
			}

			if mark == nil {
				element = cache.t2.PushFront(entry)
			} else {
				element = cache.t2.InsertAfter(entry, mark)
			}
		}

		fmt.Printf("Cache T1[%v] after:\n", cache.t1.Len())
		for e := cache.t1.Front(); e != nil; e = e.Next() {
			fmt.Printf("KEY[%v]: %v\n", e.Value.(Entry).key, e.Value.(Entry).value)
		}

		fmt.Printf("Cache T2[%v] after:\n", cache.t2.Len())
		for e := cache.t2.Front(); e != nil; e = e.Next() {
			fmt.Printf("KEY[%v]: %v\n", e.Value.(Entry).key, e.Value.(Entry).value)
		}
		fmt.Println("--------------------")

		return element.Value.(Entry).value
		// return entry.value
	} else if cache.ExistInT2(key) {
		fmt.Printf("HIT in T2\nCache T2[%v] before:\n", cache.t2.Len())
		for element = cache.t2.Front(); element != nil; element = element.Next() {
			fmt.Printf("KEY[%v]: %v\n", element.Value.(Entry).key, element.Value.(Entry).value)
		}

		cache.hitCount++
		for element = cache.t2.Front(); element != nil; element = element.Next() {
			if element.Value.(Entry).key == key {
				break
			}
		}
		// entry = element.Value.(Entry)
		// cache.t2.MoveToFront(element)
		cache.updatePosition(cache.t2, element)

		fmt.Printf("Cache T2[%v] after:\n", cache.t2.Len())
		for e := cache.t2.Front(); e != nil; e = e.Next() {
			fmt.Printf("KEY[%v]: %v\n", e.Value.(Entry).key, e.Value.(Entry).value)
		}
		fmt.Println("--------------------")
		return element.Value.(Entry).value
		// return entry.value
	}
	fmt.Println("MISS")
	fmt.Println("--------------------")
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
