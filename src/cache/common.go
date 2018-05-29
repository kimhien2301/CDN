package cache

type Storage interface {
	Len() int
	Capacity() int
	CacheList() []interface{}
	Exist(key interface{}) bool
	Insert(key, value interface{}) interface{}
	Fetch(key interface{}) interface{}
	HitCount() int
	MissCount() int
	ResetCount()
	Clear()
}

type ContentRequest struct {
	ContentKey           interface{}
	XForwardedFor        []interface{}
	TrafficWeight        float64
	TraversedServersList []string // ADD
}
