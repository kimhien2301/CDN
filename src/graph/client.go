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
		make([]string, 0), // ADD
	}
	return client.upstream.AcceptRequest(contentRequest)
}

func (client *Client_t) RandomRequest() interface{} {
	contentRequest := cache.ContentRequest{
		client.dist.Intn(),
		make([]interface{}, 0),
		client.trafficWeight,
		make([]string, 0), // ADD
	}
	return client.upstream.AcceptRequest(contentRequest)
}

func (client *Client_t) RandomRequestForInsertNewContents(count int) interface{} {
	// fmt.Println("New Content = ", count)
	requestID := client.dist.Intn()
	for requestID < count+1 {
		requestID = client.dist.Intn()
	}

	if requestID < count+1 {
		fmt.Println("Error")
	}

	contentRequest := cache.ContentRequest{
		requestID,
		make([]interface{}, 0),
		client.trafficWeight,
		make([]string, 0), // ADD
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
