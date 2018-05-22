package main

import (
	"cache"
	"parser"
)

type ServerModel_t struct {
	id             string
	storage        cache.Storage
	upstreamRouter Router
	isOrigin       bool
}

func newServer(id string, isOrigin bool) *ServerModel_t {
	server := new(ServerModel_t)
	server.id = id
	server.isOrigin = isOrigin
	return server
}

func (server *ServerModel_t) ID() string {
	return server.id
}

func (server *ServerModel_t) Storage() cache.Storage {
	return server.storage
}

func (server *ServerModel_t) setStorage(storage cache.Storage) {
	server.storage = storage
}

func (server *ServerModel_t) setUpstreamRouter(router Router) {
	server.upstreamRouter = router
}

func (server *ServerModel_t) AcceptRequest(request cache.ContentRequest) interface{} {
	// fmt.Println("FETCH 1 ")
	// fmt.Println(request.ContentKey)
	cachedData := server.Storage().Fetch(request.ContentKey)
	// fmt.Println("FETCH 2 ")
	// fmt.Println(request.ContentKey)
	if cachedData != nil {
		return cachedData
	}
	request.XForwardedFor = append(request.XForwardedFor, server.id)
	// fmt.Println("Insert 1 ")
	// fmt.Println(request.ContentKey)
	surrogateData := server.upstreamRouter.ForwardRequest(server.id, request)
	// fmt.Println("Insert 2 ")
	// fmt.Println(request.ContentKey)
	if !parser.GA {
		server.Storage().Insert(request.ContentKey, surrogateData)
	}
	// fmt.Println("Insert 3 ")
	// fmt.Println(request.ContentKey)

	// fmt.Println(parser.GA)
	return surrogateData
}

func (server *ServerModel_t) hitRate() float64 {
	hit := server.Storage().HitCount()
	miss := server.Storage().MissCount()
	return float64(hit) / (float64(hit) + float64(miss))
}
