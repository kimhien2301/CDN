package main

import "fmt"

type SeparatorRanks_t struct {
	network *Graph_t
}

func newSeparatorRanks(network *Graph_t) *SeparatorRanks_t {
	separator := new(SeparatorRanks_t)
	separator.network = network
	return separator
}

func test() {
	fmt.Println("Testing")
}
