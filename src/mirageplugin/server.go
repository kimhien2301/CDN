package main

import "cache"

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
	cachedData := server.Storage().Fetch(request.ContentKey)
	if cachedData != nil {
		return cachedData
	}
	request.XForwardedFor = append(request.XForwardedFor, server.id)
	surrogateData := server.upstreamRouter.ForwardRequest(server.id, request)
	server.Storage().Insert(request.ContentKey, surrogateData)
	return surrogateData
}

func (server *ServerModel_t) hitRate() float64 {
	hit := server.Storage().HitCount()
	miss := server.Storage().MissCount()
	return float64(hit) / (float64(hit) + float64(miss))
}
