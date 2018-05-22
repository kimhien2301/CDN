package main

import (
	"cache/eviction/iris"
	"graph"
	"math/rand"
	"utils"
)

type Node_t struct {
	id           string
	entity       interface{}
	outputLinks  []*UnidirectionalLink_t
	inputLinks   []*UnidirectionalLink_t
	dijkstraDone bool
	dijkstraCost float64
}

type Graph_t struct {
	nodes           []*Node_t
	links           []*UnidirectionalLink_t
	clients         []*graph.Client_t
	router          Router
	spectrumManager *SpectrumManager_t
	// add
	// separatorRanks *SeparatorRanks_t
	model DecodeModel
}

func (g *Graph_t) SpectrumManager() graph.SpectrumManager {
	return g.spectrumManager
}

// add
// func (g *Graph_t) SeparatorRanks() graph.SeparatorRanks {
//	  return g.separatorRanks
// }

func (g *Graph_t) LibrarySize() int {
	librarySize := 0
	for _, origin := range g.originServers() {
		storageSize := origin.Entity().(*ServerModel_t).storage.Len()
		if librarySize < storageSize {
			librarySize = storageSize
		}
	}
	return librarySize
}

func (node *Node_t) ID() string {
	return node.id
}

func (node *Node_t) Entity() interface{} {
	return node.entity
}

func (node *Node_t) OutputLinks() []graph.UnidirectionalLink {
	links := make([]graph.UnidirectionalLink, 0)
	for _, link := range node.outputLinks {
		links = append(links, graph.UnidirectionalLink(link))
	}
	return links
}

func (g *Graph_t) originServers() []*Node_t {
	originServers := make([]*Node_t, 0)
	for _, node := range g.nodes {
		if node.Entity().(*ServerModel_t).isOrigin {
			originServers = append(originServers, node)
		}
	}
	return originServers
}

func (g *Graph_t) OriginServers() []graph.Node {
	originServers := make([]graph.Node, 0)
	for _, node := range g.nodes {
		if node.Entity().(*ServerModel_t).isOrigin {
			originServers = append(originServers, graph.Node(node))
		}
	}
	return originServers
}

func (g *Graph_t) cacheServers() []*Node_t {
	cacheServers := make([]*Node_t, 0)
	for _, node := range g.nodes {
		if !node.Entity().(*ServerModel_t).isOrigin {
			cacheServers = append(cacheServers, node)
		}
	}
	return cacheServers
}

func (g *Graph_t) Links() []graph.UnidirectionalLink {
	links := make([]graph.UnidirectionalLink, 0)
	for _, link := range g.links {
		links = append(links, graph.UnidirectionalLink(link))
	}
	return links
}

func (g *Graph_t) Clients() []graph.Client {
	clients := make([]graph.Client, 0)
	for _, c := range g.clients {
		clients = append(clients, graph.Client(c))
	}
	return clients
}

func newGraph() *Graph_t {
	g := new(Graph_t)
	g.nodes = make([]*Node_t, 0)
	g.links = make([]*UnidirectionalLink_t, 0)
	g.router = NewRouter(g)
	return g
}

func newNode(id string, entity interface{}) *Node_t {
	node := new(Node_t)
	node.id = id
	node.entity = entity
	node.outputLinks = make([]*UnidirectionalLink_t, 0)
	node.inputLinks = make([]*UnidirectionalLink_t, 0)
	return node
}

// add
// func (g *Graph_t) initSeparatorRanks() {
// 	g.separatorRanks = newSeparatorRanks(g)
// }

func (g *Graph_t) initSpectrums() {
	g.spectrumManager = newSpectrumManager(4, g)
}

func (g *Graph_t) initRouter() {
	g.router.Init()
}

func (node *Node_t) opposite(link *UnidirectionalLink_t) *Node_t {
	if link.src == node {
		return link.dst
	} else if link.dst == node {
		return link.src
	}
	return nil
}

func (node *Node_t) inputAdjacent() []*Node_t {
	adjacent := make([]*Node_t, 0)
	for _, link := range node.inputLinks {
		if link.dst == node {
			adjacent = append(adjacent, link.src)
		}
	}
	return adjacent
}

func (node *Node_t) outputAdjacent() []*Node_t {
	adjacent := make([]*Node_t, 0)
	for _, link := range node.outputLinks {
		if link.src == node {
			adjacent = append(adjacent, link.dst)
		}
	}
	return adjacent
}

func (g *Graph_t) addNode(node *Node_t) {
	g.nodes = append(g.nodes, node)
}

func (g *Graph_t) detectNode(id string) *Node_t {
	for _, node := range g.nodes {
		if node.ID() == id {
			return node
		}
	}
	return nil
}

func (g *Graph_t) detectLink(src, dst *Node_t) *UnidirectionalLink_t {
	for _, link := range g.links {
		if link.src == src && link.dst == dst {
			return link
		}
	}
	return nil
}

func (g *Graph_t) connect(src, dst *Node_t, link *UnidirectionalLink_t) {
	if g.detectLink(src, dst) != nil {
		return
	}
	link.src = src
	link.dst = dst
	src.outputLinks = append(src.outputLinks, link)
	dst.inputLinks = append(dst.inputLinks, link)
	g.links = append(g.links, link)
}

func (g *Graph_t) randomRequest() interface{} {
	return g.clients[rand.Intn(len(g.clients))].RandomRequest()
}

func (g *Graph_t) ResetCounters() {
	for _, node := range g.nodes {
		node.Entity().(*ServerModel_t).Storage().ResetCount()
	}
	for _, link := range g.links {
		link.traffic = 0.0
	}
}

func (g *Graph_t) MatchStoreData(best utils.BestSeparatorRanks, bitSize int) bool {
	if best.NetworkID != g.model.NetworkID {
		return false
	}
	if len(best.Ranks) != bitSize {
		return false
	}
	if best.RequestModel != g.model.Clients[0].RequestModelID {
		return false
	}
	for _, requestModel := range g.model.RequestModels {
		if requestModel.ID == best.RequestModel {
			for index := range best.RequestModelParams {
				if best.RequestModelParams[index] != requestModel.Parameters[index] {
					return false
				}
			}
		}
	}
	return true
}

func (g *Graph_t) GenerateBestSeparatorRanksData(bestSeparatorRanks []int) utils.BestSeparatorRanks {
	var best utils.BestSeparatorRanks
	best.NetworkID = g.model.NetworkID
	best.RequestModel = g.model.Clients[0].RequestModelID
	for _, requestModel := range g.model.RequestModels {
		if requestModel.ID == best.RequestModel {
			best.RequestModelParams = requestModel.Parameters
			break
		}
	}
	best.Ranks = bestSeparatorRanks
	return best
}

func (g *Graph_t) CacheAlgorithm() string {
	return g.model.Nodes[0].CacheAlgorithm
}

func (g *Graph_t) SpectrumCapacity() int {
	return g.cacheServers()[0].Entity().(graph.ServerModel).Storage().(iris.Accessor).SpectrumCapacity()
}

// ADD
func (g *Graph_t) GetExpectTraffic() float64 {
	return g.expectTrafficWithNoFillCaches().totalTraffic()
}

func (g *Graph_t) SetCacheServers(chromosome [][]int) {
	for i, client := range g.clients {
		cacheServer := client.Upstream().Storage()
		cacheServer.Clear()
		for j := 0; j < cacheServer.Capacity(); j++ {
			cacheServer.Insert(chromosome[i][j], chromosome[i][j])
		}
	}
}

func (g *Graph_t) GetNumberOfCacheServers() int {
	return len(g.clients)
}

func (g *Graph_t) GetCacheCapacity() int {
	return g.clients[0].Upstream().Storage().Capacity()
}
