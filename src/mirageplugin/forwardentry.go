package main

import "graph"

type ForwardEntry_t struct {
	node graph.Node
	link graph.UnidirectionalLink
	cost float64
}

type ForwardEntry interface {
	Node() graph.Node
	Link() graph.UnidirectionalLink
	Cost() float64
}

func (entry *ForwardEntry_t) Node() graph.Node {
	return entry.node
}

func (entry *ForwardEntry_t) Link() graph.UnidirectionalLink {
	return entry.link
}

func (entry *ForwardEntry_t) Cost() float64 {
	return entry.cost
}
