package main

import "fmt"
import "parser"
import "cache"
import "utils"
import "math"

type Router_t struct {
	graph        *Graph_t
	routingTable map[string]map[string][]ForwardEntry
}

type Router interface {
	selectForwardEntryBySpectrum(string, cache.ContentRequest) ForwardEntry
	selectForwardEntry(string, cache.ContentRequest) ForwardEntry
	Init()
	ForwardRequest(string, cache.ContentRequest) interface{}
	RoutingTable() map[string]map[string][]ForwardEntry
	// ADD
	inspectAllRoutingTables()
	ForwardRequestForIris(string, cache.ContentRequest) interface{}
}

func NewRouter(g *Graph_t) Router {
	router := new(Router_t)
	router.graph = g
	return router
}

func (router *Router_t) RoutingTable() map[string]map[string][]ForwardEntry {
	return router.routingTable
}

func (router *Router_t) dijkstra(dstNode *Node_t) {
	router.resetDijkstraVariables()
	dstNode.dijkstraCost = 0.0
	var doneNode *Node_t

	for {
		doneNode = nil
		for _, node := range router.graph.nodes {
			if node.dijkstraDone || node.dijkstraCost < 0 {
				continue
			}
			if doneNode == nil || node.dijkstraCost < doneNode.dijkstraCost {
				doneNode = node
			}
		}
		if doneNode == nil {
			break
		}

		doneNode.dijkstraDone = true
		for i := 0; i < len(doneNode.outputLinks); i++ {
			to := doneNode.outputLinks[i].dst
			cost := doneNode.dijkstraCost + doneNode.outputLinks[i].cost
			if to.dijkstraCost < 0 || cost < to.dijkstraCost {
				to.dijkstraCost = cost
			}
		}

	}
}

func (router *Router_t) createRoutingTable() {
	router.routingTable = make(map[string]map[string][]ForwardEntry)

	for _, dstNode := range router.graph.nodes {
		router.dijkstra(dstNode)

		for _, srcNode := range router.graph.nodes {
			if router.routingTable[srcNode.ID()] == nil {
				router.routingTable[srcNode.ID()] = make(map[string][]ForwardEntry)
			}
			if router.routingTable[srcNode.ID()][dstNode.ID()] == nil {
				router.routingTable[srcNode.ID()][dstNode.ID()] = make([]ForwardEntry, 0)
			}
			if srcNode == dstNode {
				continue
			}
			minCost := math.Inf(0)

			for _, link := range srcNode.inputLinks {
				if minCost > srcNode.opposite(link).dijkstraCost {
					minCost = srcNode.opposite(link).dijkstraCost
				}
			}
			for _, link := range srcNode.inputLinks {
				if srcNode.opposite(link).dijkstraCost == minCost {
					forwardEntry := new(ForwardEntry_t)
					forwardEntry.node = srcNode.opposite(link)
					forwardEntry.link = link
					forwardEntry.cost = srcNode.dijkstraCost
					router.routingTable[srcNode.ID()][dstNode.ID()] =
						append(router.routingTable[srcNode.ID()][dstNode.ID()], forwardEntry)
				}
			}
		}

	}
}

func (router *Router_t) resetDijkstraVariables() {
	for _, node := range router.graph.nodes {
		node.dijkstraDone = false
		node.dijkstraCost = -1.0
	}
}

func (router *Router_t) inspectForwardEntry(entry ForwardEntry) {
	utils.DebugPrint(fmt.Sprintf("    <next: %s, via_link: %s -> %s, cost:%f>\n", entry.Node().ID(), entry.Link().Src().ID(), entry.Link().Dst().ID(), entry.Cost()))
}

func (router *Router_t) inspectRoutingTable(fromNode *Node_t) {
	utils.DebugPrint(fmt.Sprintf("From: %s\n", fromNode.ID()))
	for _, destNode := range router.graph.nodes {
		utils.DebugPrint(fmt.Sprintf("  Dest: %s\n", destNode.ID()))
		for _, forwardEntry := range router.routingTable[fromNode.ID()][destNode.ID()] {
			router.inspectForwardEntry(forwardEntry)
		}
	}
}

func (router *Router_t) inspectAllRoutingTables() {
	for _, fromNode := range router.graph.nodes {
		router.inspectRoutingTable(fromNode)
	}
}

func (router *Router_t) Init() {
	router.createRoutingTable()
}

func (router *Router_t) selectForwardEntryBySpectrum(fromID string, request cache.ContentRequest) ForwardEntry {
	// ADD
	srcNode := router.graph.detectNode(fromID)
	contentTags := router.graph.spectrumManager.spectrumTags[request.ContentKey.(int)]
	requestRoutingTable := router.graph.spectrumManager.spectrumRoutingTable[srcNode]

	var forwardEntry ForwardEntry
	if forwardEntry != nil {
		fmt.Println("Error: forwardEntry in selectForwardEntryBySpectrum")
	}

	// check if content has a color tag
	if len(contentTags) > 0 {
		for colorMask, forwardEntries := range requestRoutingTable {
			if len(forwardEntries) > 1 {
				fmt.Println("Error: len(forwardEntries) in selectForwardEntryBySpectrum")
				// fmt.Println(requestRoutingTable)
				// fmt.Println(forwardEntries)
			}
			// fmt.Println(contentTags)
			// fmt.Println(colorMask)
			// fmt.Println(forwardEntries)
			if contentTags[colorMask] == 1 && !inTraversedServersList(forwardEntries[0].Node().ID(), request) {
				if forwardEntry == nil {
					forwardEntry = forwardEntries[0]
				} else {
					if forwardEntries[0].Cost() < forwardEntry.Cost() {
						forwardEntry = forwardEntries[0]
					}
				}
			}

		}
	}

	if forwardEntry == nil {
		forwardEntry = router.selectForwardEntry(fromID, request)
	}
	return forwardEntry
}

func inTraversedServersList(nodeID string, request cache.ContentRequest) bool {
	for i := range request.TraversedServersList {
		if nodeID == request.TraversedServersList[i] {
			return true
		}
	}
	return false
}

func (router *Router_t) selectForwardEntry(fromID string, request cache.ContentRequest) ForwardEntry {
	srcNode := router.graph.detectNode(fromID)
	originServers := router.graph.originServers()
	forwardEntries := make([]ForwardEntry, 0)

	for _, origin := range originServers {
		for _, forwardEntry := range router.routingTable[srcNode.ID()][origin.ID()] {
			forwardEntries = append(forwardEntries, forwardEntry)
		}
	}

	minCost := math.Inf(0)
	for _, forwardEntry := range forwardEntries {
		if minCost > forwardEntry.Cost() {
			minCost = forwardEntry.Cost()
		}
	}

	minCostEntries := make([]ForwardEntry, 0)
	for _, forwardEntry := range forwardEntries {
		if forwardEntry.Cost() == minCost {
			minCostEntries = append(minCostEntries, forwardEntry)
		}
	}

	if len(minCostEntries) > 0 {
		return minCostEntries[0]
	}
	return nil
}

func (router *Router_t) ForwardRequest(fromID string, request cache.ContentRequest) interface{} {
	var forwardEntry ForwardEntry
	if router.graph.model.Nodes[0].CacheAlgorithm == "iris" && !parser.Options.UseShortestPath {
		// forwardEntry = router.selectForwardEntryBySpectrum(fromID, request)
	} else {
		forwardEntry = router.selectForwardEntry(fromID, request)
	}
	surrogate := forwardEntry.Node().Entity().(*ServerModel_t)
	forwardEntry.Link().(*UnidirectionalLink_t).traffic += request.TrafficWeight
	return surrogate.AcceptRequest(request)
}

// ADD
func (router *Router_t) ForwardRequestForIris(fromID string, request cache.ContentRequest) interface{} {
	var forwardEntry ForwardEntry
	var surrogate *ServerModel_t
	var surrogateData interface{}

	// Requesting data
	//request.TraversedServersList = make([]string, 0)
	request.TraversedServersList = append(request.TraversedServersList, fromID)

	forwardEntry = router.selectForwardEntryBySpectrum(fromID, request)
	// fmt.Println("Fisrt forwardEntry from node")
	// fmt.Println(fromID)

	for i := 0; surrogateData == nil; i++ {
		surrogate = forwardEntry.Node().Entity().(*ServerModel_t)
		// forwardEntry.Link().(*UnidirectionalLink_t).traffic += request.TrafficWeight
		surrogateData = surrogate.Storage().Fetch(request.ContentKey)

		request.TraversedServersList = append(request.TraversedServersList, surrogate.ID())

		// check node already traversed, default interface in EntryBySpectrum
		forwardEntry = router.selectForwardEntryBySpectrum(surrogate.ID(), request)

		if i < 3 {
			// fmt.Println("forward to find data")
			// fmt.Println(surrogate.ID())
			// fmt.Println(surrogateData)
			// fmt.Println("TraversedServersList")
			// fmt.Println(request.TraversedServersList)
		}
	}

	// Reurning data
	forwardEntry = router.selectForwardEntryToNode(surrogate.ID(), fromID)
	// fmt.Println("Returning data")
	// fmt.Println("Server has the content")
	// fmt.Println(surrogate.ID())

	for surrogate.ID() != fromID {
		// fmt.Println("server cached")
		// fmt.Println(surrogate.ID())

		//forwardEntry.Link().(*UnidirectionalLink_t).traffic += request.TrafficWeight

		router.selectForwardEntryToNode(forwardEntry.Node().ID(), surrogate.ID()).Link().(*UnidirectionalLink_t).traffic += request.TrafficWeight

		// cache content on cache servers which are on the path to first requested server
		surrogate.Storage().Insert(request.ContentKey, surrogateData)

		surrogate = forwardEntry.Node().Entity().(*ServerModel_t)
		forwardEntry = router.selectForwardEntryToNode(surrogate.ID(), fromID)
	}

	return surrogateData
}

func (router *Router_t) selectForwardEntryToNode(fromID, toID string) ForwardEntry {
	forwardEntries := router.routingTable[fromID][toID]
	var forwardEntry ForwardEntry
	minCost := math.Inf(0)

	for i := range forwardEntries {
		if minCost > forwardEntries[i].Cost() {
			minCost = forwardEntries[i].Cost()
			forwardEntry = forwardEntries[i]
		}
	}

	return forwardEntry
}
