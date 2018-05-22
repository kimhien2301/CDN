package origin

type Library struct {
	size int
	hit  int
}

func NewLibrary(size int) *Library {
	// fmt.Printf("New library[%v]\n", size)
	library := new(Library)
	library.size = size
	return library
}

func (library *Library) CacheList() []interface{} {
	cacheList := make([]interface{}, 0)
	for id := 0; id < library.size; id++ {
		cacheList = append(cacheList, id+1)
	}
	return cacheList
}

func (library *Library) HitCount() int {
	return library.hit
}

func (library *Library) MissCount() int {
	return 0
}

func (library *Library) ResetCount() {
	library.hit = 0
}

func (library *Library) Capacity() int {
	return library.size
}

func (library *Library) Len() int {
	return library.size
}

func (library *Library) Exist(key interface{}) bool {
	// fmt.Println("TRUE")
	return true
}

func (library *Library) Insert(key, value interface{}) interface{} {
	return value
}

func (library *Library) Fetch(key interface{}) interface{} {
	// fmt.Println("HIT")
	library.hit++
	return key
}

func (library *Library) Clear() {
	library.ResetCount()
}
