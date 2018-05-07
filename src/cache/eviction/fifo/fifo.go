package fifo

import "container/list"
import "sync"

type Entry struct {
	key, value interface{}
}

type FIFO struct {
	length    int
	list      *list.List
	hit, miss int
	mutex     sync.Mutex
}

func New(length int) *FIFO {
	fifo := new(FIFO)
	fifo.length = length
	fifo.list = list.New()
	fifo.hit = 0
	fifo.miss = 0
	return fifo
}

func (fifo *FIFO) CacheList() []interface{} {
	cacheList := make([]interface{}, 0)
	for element := fifo.list.Front(); element != nil; element = element.Next() {
		cacheList = append(cacheList, element.Value.(Entry).key)
	}
	return cacheList
}

func (fifo *FIFO) Clear() {
	fifo.list = list.New()
	fifo.ResetCount()
}

func (fifo *FIFO) Len() int {
	return fifo.list.Len()
}

func (fifo *FIFO) Exist(key interface{}) bool {
	for element := fifo.list.Front(); element != nil; element = element.Next() {
		if element == nil {
			break
		}
		if element.Value.(Entry).key == key {
			return true
		}
	}
	return false
}

func (fifo *FIFO) Insert(key, value interface{}) interface{} {
	fifo.mutex.Lock()
	fifo.list.PushFront(Entry{key, value})
	if fifo.Len() > fifo.length {
		fifo.list.Remove(fifo.list.Back())
	}
	fifo.mutex.Unlock()
	return value
}

func (fifo *FIFO) Fetch(key interface{}) interface{} {
	for element := fifo.list.Front(); element != nil; element = element.Next() {
		if element == nil {
			break
		}
		if element.Value.(Entry).key == key {
			fifo.hit++
			return element.Value.(Entry).value
		}
	}
	fifo.miss++
	return nil
}

func (fifo *FIFO) Capacity() int {
	return fifo.length
}

func (fifo *FIFO) HitCount() int {
	return fifo.hit
}

func (fifo *FIFO) MissCount() int {
	return fifo.miss
}

func (fifo *FIFO) ResetCount() {
	fifo.hit = 0
	fifo.miss = 0
}
