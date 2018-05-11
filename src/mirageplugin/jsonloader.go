package main

import (
	"cache"
	"cache/eviction/admission"
	"cache/eviction/arc"
	"cache/eviction/fifo"
	"cache/eviction/iclfu"
	"cache/eviction/iris"
	"cache/eviction/lfu"
	"cache/eviction/lirs"
	"cache/eviction/lru"
	"cache/eviction/lruk"
	"cache/eviction/modifiedlru"
	"cache/eviction/nocache"
	"cache/eviction/random"
	"cache/eviction/srrip"
	"cache/eviction/wlfu"
	"distribution"
	"distribution/gamma"
	"distribution/userdist"
	"distribution/zipf"
	"encoding/json"
	"fmt"
	"graph"
	"io/ioutil"
	"origin"
)

type RequestDecodeModel struct {
	ID            string    `json:"Id"`
	Model         string    `json:"Model"`
	ParameterKeys []string  `json:"ParameterKeys"`
	Parameters    []float64 `json:"Parameters"`
}

type OriginDecodeModel struct {
	ID          string `json:"Id"`
	LibrarySize int    `json:"LibrarySize"`
}

type NodeDecodeModel struct {
	ID             string    `json:"Id"`
	CacheAlgorithm string    `json:"CacheAlgorithm"`
	ParameterKeys  []string  `json:"ParameterKeys"`
	Parameters     []float64 `json:"Parameters"`
}

type LinkDecodeModel struct {
	EdgeNodeIds   []string `json:"EdgeNodeIds"`
	Cost          float64  `json:"Cost"`
	Bidirectional bool     `json:"Bidirectional"`
}

type ClientDecodeModel struct {
	RequestModelID string `json:"RequestModelId"`
	UpstreamID     string `json:"UpstreamId"`
	TrafficWeight  float64
}

type DecodeModel struct {
	NetworkID      string               `json:"NetworkId"`
	RequestModels  []RequestDecodeModel `json:"RequestModels"`
	Origin         OriginDecodeModel    `json:"OriginServer"`
	Nodes          []NodeDecodeModel    `json:"CacheServers"`
	Links          []LinkDecodeModel    `json:"Links"`
	Clients        []ClientDecodeModel  `json:"Clients"`
	ModelRelations map[interface{}]interface{}
}

func loadFile(filename string) DecodeModel {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("Fatal error in reading json file: %s", err))
	}

	var graphDecodeModel DecodeModel
	if err = json.Unmarshal(bytes, &graphDecodeModel); err != nil {
		panic(fmt.Errorf("Fatal error in Unmarshal(): %s", err))
	}

	return graphDecodeModel
}

func (g *Graph_t) loadOriginServer(graphDecodeModel DecodeModel) {
	originServer := newServer(graphDecodeModel.Origin.ID, true)
	library := origin.NewLibrary(int(graphDecodeModel.Origin.LibrarySize))
	originServer.setStorage(library)
	originNode := newNode(graphDecodeModel.Origin.ID, originServer)
	g.addNode(originNode)
	graphDecodeModel.ModelRelations[originNode] = graphDecodeModel.Origin
}

func (g *Graph_t) loadCacheServers(graphDecodeModel DecodeModel) {
	for _, nodeModel := range graphDecodeModel.Nodes {

		params := make(map[string]float64)
		for index, key := range nodeModel.ParameterKeys {
			params[key] = nodeModel.Parameters[index]
		}

		edgeServer := newServer(nodeModel.ID, false)

		var storage cache.Storage
		switch nodeModel.CacheAlgorithm {
		case "modifiedlru":
			storage = modifiedlru.New(int(params["Capacity"]), int(params["Jump"]))
		case "random":
			storage = random.New(int(params["Capacity"]))
		case "lru":
			storage = lru.New(int(params["Capacity"]))
		case "lruk":
			storage = lruk.New(int(params["Capacity"]), int(params["K"]))
		case "srrip":
			storage = srrip.New(int(params["Capacity"]), int(params["RRPVbit"]))
		case "lirs":
			storage = lirs.New(int(params["Capacity"]))
		case "arc":
			storage = arc.New(int(params["Capacity"]))
		case "fifo":
			storage = fifo.New(int(params["Capacity"]))
		case "iclfu":
			storage = iclfu.New(int(params["Capacity"]))
		case "windowlfu":
			storage = windowlfu.New(int(params["Capacity"]), int(params["Window"]))
		case "lfu":
			storage = lfu.New(int(params["Capacity"]))
		case "admission":
			admissionList := make([]interface{}, 0)
			for rank := 0; rank < int(params["Capacity"]); rank++ {
				admissionList = append(admissionList, rank)
			}
			storage = admission.New(admissionList)
		case "iris":
			storage = NewIrisCache(int(params["Capacity"]), params["SpectrumRatio"])
		case "nocache":
			storage = nocache.New(int(params["Capacity"]))
		}
		edgeServer.setStorage(storage)
		edgeServer.setUpstreamRouter(g.router)

		cacheNode := newNode(nodeModel.ID, edgeServer)
		g.addNode(cacheNode)
		graphDecodeModel.ModelRelations[cacheNode] = nodeModel
	}
}

func (g *Graph_t) loadLinks(graphDecodeModel DecodeModel) {
	for _, linkModel := range graphDecodeModel.Links {
		SrcID := linkModel.EdgeNodeIds[0]
		DstID := linkModel.EdgeNodeIds[1]

		forwardLink := new(UnidirectionalLink_t)
		forwardLink.cost = linkModel.Cost
		g.connect(g.detectNode(SrcID), g.detectNode(DstID), forwardLink)
		graphDecodeModel.ModelRelations[forwardLink] = linkModel
		if linkModel.Bidirectional {
			inverseLink := new(UnidirectionalLink_t)
			inverseLink.cost = linkModel.Cost
			g.connect(g.detectNode(DstID), g.detectNode(SrcID), inverseLink)
			graphDecodeModel.ModelRelations[inverseLink] = linkModel
		}
	}
}

func loadDistributions(graphDecodeModel DecodeModel, librarySize int) map[string]distribution.Distribution {
	dists := make(map[string]distribution.Distribution)
	for _, requestModel := range graphDecodeModel.RequestModels {
		params := make(map[string]float64)
		for index, key := range requestModel.ParameterKeys {
			params[key] = requestModel.Parameters[index]
		}
		switch requestModel.Model {
		case "gamma":
			dists[requestModel.ID] = gamma.New(params["K"], params["Theta"], librarySize)
		case "zipf":
			dists[requestModel.ID] = zipf.New(params["Skewness"], librarySize)
		case "userdist":
			dists[requestModel.ID] = userdist.New("id_count_full.csv")
		}
	}
	return dists
}

func (g *Graph_t) loadClients(graphDecodeModel DecodeModel) {
	librarySize := g.LibrarySize()
	dists := loadDistributions(graphDecodeModel, librarySize)
	for _, clientModel := range graphDecodeModel.Clients {
		cl := graph.NewClient(dists[clientModel.RequestModelID], g.detectNode(clientModel.UpstreamID).Entity().(*ServerModel_t), clientModel.TrafficWeight)
		g.clients = append(g.clients, cl)
		graphDecodeModel.ModelRelations[cl] = clientModel
	}
}

func (g *Graph_t) clone() *Graph_t {
	return loadModel(g.model)
}

func loadModel(graphDecodeModel DecodeModel) *Graph_t {
	graphDecodeModel.ModelRelations = make(map[interface{}]interface{})

	network := newGraph()
	network.model = graphDecodeModel
	network.loadOriginServer(graphDecodeModel)
	network.loadCacheServers(graphDecodeModel)
	network.loadLinks(graphDecodeModel)
	network.loadClients(graphDecodeModel)

	network.initRouter()

	if graphDecodeModel.Nodes[0].CacheAlgorithm == "iris" {
		network.initSpectrums()
		network.SpectrumManager().(*SpectrumManager_t).setServerSpectrums()

		// set server spectrum to cache.go
		for _, node := range network.cacheServers() {
			node.Entity().(graph.ServerModel).Storage().(iris.Accessor).SetServerSpectrum(network.spectrumManager.serverSpectrums[node])
		}
	}

	// add
	// network.initSpectrums()
	// network.SpectrumManager().(*SpectrumManager_t).setServerSpectrums()
	// network.initSeparatorRanks()

	return network
}

func LoadGraph(filename string) graph.Graph {
	graphDecodeModel := loadFile(filename)
	return loadModel(graphDecodeModel)
}
