package graph

import (
	"cache"
	"distribution"
	"fmt"
)

type Client_t struct {
	dist          distribution.Distribution
	upstream      ServerModel
	trafficWeight float64
}

func NewClient(dist distribution.Distribution, upstream ServerModel, trafficWeight float64) *Client_t {
	client := new(Client_t)
	client.dist = dist
	client.upstream = upstream
	client.trafficWeight = trafficWeight
	return client
}

func (client *Client_t) RequestByID(contentID int) interface{} {
	contentRequest := cache.ContentRequest{
		contentID,
		make([]interface{}, 0),
		client.trafficWeight,
	}

	return client.upstream.AcceptRequest(contentRequest)
}

func (client *Client_t) RandomRequest() interface{} {
	contentRequest := cache.ContentRequest{
		client.dist.Intn(),
		make([]interface{}, 0),
		client.trafficWeight,
	}
	fmt.Printf("Content request: %v.\n", contentRequest.ContentKey)
	return client.upstream.AcceptRequest(contentRequest)
}

func (client *Client_t) RandomRequestForInsertNewContents() interface{} {
	requestID := client.dist.Intn()
	for requestID < 6 {
		requestID = client.dist.Intn()
	}

	contentRequest := cache.ContentRequest{
		requestID,
		make([]interface{}, 0),
		client.trafficWeight,
	}
	// fmt.Printf("Content request: %v.\n", contentRequest.ContentKey)
	return client.upstream.AcceptRequest(contentRequest)
}

func (client *Client_t) Upstream() ServerModel {
	return client.upstream
}

func (client *Client_t) Dist() distribution.Distribution {
	return client.dist
}

func (client *Client_t) TrafficWeight() float64 {
	return client.trafficWeight
}
