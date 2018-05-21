package arcz

import (
	"container/list"
	"math"
	"time"
)

type CacheStorage struct {
	t1        *list.List
	b1        *list.List
	t2        *list.List
	b2        *list.List
	Z         *list.List
	p         int
	k         int
	capacity  int
	hitCount  int
	missCount int
}

type Entry struct {
	key        interface{}
	value      interface{}
	accessTime time.Time
}

func New(capacity, k int) *CacheStorage {
	cache := new(CacheStorage)
	cache.t1 = list.New()
	cache.b1 = list.New()
	cache.t2 = list.New()
	cache.b2 = list.New()
	cache.Z = list.New()
	cache.k = k
	cache.capacity = capacity
	cache.p = 0
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

func (cache *CacheStorage) ExistInZ(key interface{}) bool {
	for element := cache.Z.Front(); element != nil; element = element.Next() {
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
		cache.t2.PushFront(Entry{key, value, time.Now()})
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
		cache.t2.PushFront(Entry{key, value, time.Now()})
	} else if !cache.ExistInAny(key) {
		if cache.t1.Len()+cache.b1.Len() == cache.capacity {
			var lruPage *list.Element
			if cache.t1.Len() < cache.capacity {
				lruPage = cache.b1.Back()
				cache.b1.Remove(lruPage)
				cache.Replace(key)
			} else {
				lruPage = cache.t1.Back()
				cache.t1.Remove(lruPage)
			}
			cache.Z.PushFront(lruPage.Value.(Entry))
			if cache.Z.Len() > cache.k {
				cache.Z.Remove(cache.Z.Back())
			}
		} else if cache.t1.Len()+cache.b1.Len() < cache.capacity {
			if cache.t1.Len()+cache.t2.Len()+cache.b1.Len()+cache.b2.Len() >= cache.capacity {
				if cache.t1.Len()+cache.t2.Len()+cache.b1.Len()+cache.b2.Len() == 2*cache.capacity {
					// fmt.Printf("Old LRU of B2: KEY[%v]: %v\n", cache.b2.Back().Value.(Entry).key, cache.b2.Back().Value.(Entry).accessTime)
					cache.b2.Remove(cache.b2.Back())
				}
				cache.Replace(key)

				if cache.b2.Back() != nil && cache.Z.Back() != nil {
					// fmt.Printf("LRU of Z: KEY[%v]: %v\n", cache.Z.Back().Value.(Entry).key, cache.Z.Back().Value.(Entry).accessTime)
					deleteCandidates := make([]*list.Element, 0)
					for element := cache.Z.Back(); element != nil; element = element.Prev() {
						if element.Value.(Entry).accessTime.Before(cache.b2.Back().Value.(Entry).accessTime) {
							deleteCandidates = append(deleteCandidates, element)
							// fmt.Printf("Element: KEY[%v] %v %v\n", element.Value.(Entry).key, element.Value.(Entry).value, element.Value.(Entry).accessTime)
						}
					}
					// fmt.Printf("List Z[%d]: \n", cache.Z.Len())
					// for element := cache.Z.Front(); element != nil; element = element.Next() {
					// 	fmt.Printf("KEY: %v VALUE: %v TIME: %v %v\n", element.Value.(Entry).key, element.Value.(Entry).value, element.Value.(Entry).accessTime, reflect.TypeOf(element.Value.(Entry).accessTime))
					// }
					// fmt.Printf("New LRU of B2: KEY[%v]: %v %v\n", cache.b2.Back().Value.(Entry).key, cache.b2.Back().Value.(Entry).accessTime, reflect.TypeOf(cache.b2.Back()))
					// fmt.Printf("Delete candidates[%v]: \n", len(deleteCandidates))
					// for _, element := range deleteCandidates {
					// 	fmt.Printf("KEY[%v]: %v %v\n", element.Value.(Entry).key, element.Value.(Entry).value, element.Value.(Entry).accessTime)
					// }
					for _, element := range deleteCandidates {
						cache.Z.Remove(element)
					}
					// fmt.Printf("List Z[%d] after: \n", cache.Z.Len())
					// for element := cache.Z.Front(); element != nil; element = element.Next() {
					// 	fmt.Printf("KEY: %v VALUE: %v TIME: %v %v\n", element.Value.(Entry).key, element.Value.(Entry).value, element.Value.(Entry).accessTime, reflect.TypeOf(element.Value.(Entry).accessTime))
					// }
				}
			}
		}
		cache.t1.PushFront(Entry{key, value, time.Now()})
	}
	return value
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	var element *list.Element
	var entry Entry

	if cache.ExistInT1(key) {
		// fmt.Println("HIT in T1")
		cache.hitCount++
		for element = cache.t1.Front(); element != nil; element = element.Next() {
			if element.Value.(Entry).key == key {
				break
			}
		}
		entry = element.Value.(Entry)
		entry.accessTime = time.Now()
		cache.t1.Remove(element)
		cache.t2.PushFront(entry)
		return entry.value
	} else if cache.ExistInT2(key) {
		// fmt.Println("HIT in T2")
		cache.hitCount++
		for element = cache.t2.Front(); element != nil; element = element.Next() {
			if element.Value.(Entry).key == key {
				break
			}
		}
		entry = element.Value.(Entry)
		entry.accessTime = time.Now()
		cache.t2.Remove(element)
		cache.t2.PushFront(entry)
		return entry.value
	} else if cache.ExistInZ(key) {
		// fmt.Println("HIT in Z")
		// fmt.Printf("T1[%v] T2[%v] before: \n", cache.t1.Len(), cache.t2.Len())
		// for element := cache.t2.Front(); element != nil; element = element.Next() {
		// 	fmt.Printf("KEY[%v] VALUE %v TIME %v\n", element.Value.(Entry).key, element.Value.(Entry).value, element.Value.(Entry).accessTime)
		// }
		// fmt.Printf("List Z[%d]: \n", cache.Z.Len())
		// for element := cache.Z.Front(); element != nil; element = element.Next() {
		// 	fmt.Printf("KEY: %v VALUE: %v TIME: %v %v\n", element.Value.(Entry).key, element.Value.(Entry).value, element.Value.(Entry).accessTime, reflect.TypeOf(element.Value.(Entry).accessTime))
		// }

		cache.hitCount++
		for element = cache.Z.Front(); element != nil; element = element.Next() {
			if element.Value.(Entry).key == key {
				break
			}
		}
		entry = element.Value.(Entry)
		entry.accessTime = time.Now()
		cache.Replace(key)
		cache.Z.Remove(element)
		cache.t2.PushFront(entry)

		// fmt.Printf("T1[%v] T2[%v] after: \n", cache.t1.Len(), cache.t2.Len())
		// for element := cache.t2.Front(); element != nil; element = element.Next() {
		// 	fmt.Printf("KEY[%v] VALUE %v TIME %v\n", element.Value.(Entry).key, element.Value.(Entry).value, element.Value.(Entry).accessTime)
		// }
		// fmt.Printf("List Z[%d] after HIT: \n", cache.Z.Len())
		// for element := cache.Z.Front(); element != nil; element = element.Next() {
		// 	fmt.Printf("KEY: %v VALUE: %v TIME: %v\n", element.Value.(Entry).key, element.Value.(Entry).value, element.Value.(Entry).accessTime)
		// }

		return entry.value
	}
	// fmt.Println("MISS")
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
	cache.Z = list.New()
}
