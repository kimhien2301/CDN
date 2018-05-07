package iris

type Accessor interface {
	CacheList() []interface{}
	Capacity() int
	Clear()
	Exist(interface{}) bool
	Fetch(interface{}) interface{}
	FillUp()
	HitCount() int
	Insert(interface{}, interface{}) interface{}
	Inspect()
	Len() int
	MissCount() int
	ResetCount()
	SetContentSpectrums([]uint64)
	SetServerSpectrum(uint64)
	SpectrumCapacity() int
}
