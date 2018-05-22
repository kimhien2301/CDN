package graph

import (
	"cache"
	"distribution"
	"utils"
)

type Graph interface {
	PlainReport()
	JsonReport()
	CacheAlgorithm() string
	LibrarySize() int
	SpectrumManager() SpectrumManager
	MatchStoreData(utils.BestSeparatorRanks, int) bool
	GenerateBestSeparatorRanksData([]int) utils.BestSeparatorRanks
	Clients() []Client
	ResetCounters()
	Links() []UnidirectionalLink
	OriginServers() []Node
	SpectrumCapacity() int

	// ADD
	GetExpectTraffic() float64
	SetCacheServers(chromosome [][]int)
	GetNumberOfCacheServers() int
	GetCacheCapacity() int
	GetCacheHitRate() float64
}

type Client interface {
	RandomRequest() interface{}
	RequestByID(int) interface{}
	Upstream() ServerModel
	// ADD
	Dist() distribution.Distribution
	RandomRequestForInsertNewContents() interface{}
}

type SpectrumManager interface {
	BestSeparatorRanks([]int) []int
	BestReferenceRanks(utils.MirageStore) []int
	SetContentSpectrums([]int)
}

// ADD
type SeparatorRanks interface {
	GetSeparatorRanks()
}

type ServerModel interface {
	Storage() cache.Storage
	AcceptRequest(cache.ContentRequest) interface{}
	ID() string
}

type UnidirectionalLink interface {
	Src() Node
	Dst() Node
	Traffic() float64
	SetTraffic(float64)
}

type Node interface {
	Entity() interface{}
	ID() string
	OutputLinks() []UnidirectionalLink
}
