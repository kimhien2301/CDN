package main

import "graph"

type UnidirectionalLink_t struct {
	src     *Node_t
	dst     *Node_t
	cost    float64
	traffic float64
}

func (link *UnidirectionalLink_t) SetTraffic(traffic float64) {
	link.traffic += traffic
}

func (link *UnidirectionalLink_t) Traffic() float64 {
	return link.traffic
}

func (link *UnidirectionalLink_t) Src() graph.Node {
	return link.src
}

func (link *UnidirectionalLink_t) Dst() graph.Node {
	return link.dst
}
