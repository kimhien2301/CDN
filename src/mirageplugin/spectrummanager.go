package main

import (
	"container/list"
	"fmt"
	"math"
	"sort"
	"utils"
)

type EstimationResult struct {
	separatorRanks []int
	trafficSize    float64
}

type SpectrumManager_t struct {
	network              *Graph_t
	baseSpectrums        []interface{}                         // colors not sure???
	spectrumTags         map[int][]uint64                      // contents' color
	bitSize              int                                   // # colors
	contentSpectrums     []uint64                              // separatorRanks
	serverSpectrums      map[*Node_t]uint64                    // servers' color
	spectrumRoutingTable map[*Node_t]map[uint64][]ForwardEntry // color-based routing table
}

func newSpectrumManager(bitSize int, network *Graph_t) *SpectrumManager_t {
	manager := new(SpectrumManager_t)
	// add
	manager.network = network
	manager.bitSize = bitSize
	manager.serverSpectrums = make(map[*Node_t]uint64)
	manager.baseSpectrums = make([]interface{}, 0)
	manager.contentSpectrums = make([]uint64, 0)
	for i := 0; i < bitSize; i++ {
		manager.baseSpectrums = append(manager.baseSpectrums, uint64(i))
	}
	return manager
}

func (manager *SpectrumManager_t) initSpectrumRoutingTable() {
}

func (manager *SpectrumManager_t) inspectSpectrumRoutingTable() {
}

func (manager *SpectrumManager_t) initSpectrums() {
}

func (manager *SpectrumManager_t) initBaseSpectrums() {
}

func (manager *SpectrumManager_t) estimateTotalTraffic(separatorRanks []int, join chan EstimationResult) {
}

func (manager *SpectrumManager_t) adjustSeparatorRanks(separatorRanks []int, librarySize int) bool {
	return false
}

func (manager *SpectrumManager_t) separatorRanksID(separatorRanks []int) string {
	return ""
}

func (manager *SpectrumManager_t) BestSeparatorRanks(separatorRanks []int) []int {
	// N := manager.bitSize                                            // # colors
	// C := manager.network.clients[0].Upstream().Storage().Capacity() // cache server capacity
	// numUsrs := len(manager.network.clients)
	// numServers := len(manager.network.nodes)
	// numContents := manager.network.LibrarySize()

	// // fmt.Printf("\nNum of Usrs: %v, Num of Servers: %v, Num of Contents: %v\n", numUsrs, numServers, numContents)

	// S := make([]int, 4)
	// S_prev := make([]int, 4)
	// var S_tmp []int
	// T_min := math.MaxFloat64

	// //fmt.Println(reflect.TypeOf(T_min))

	// S[N-1] = N * C
	// fmt.Println("Calculate separator ranks:")
	// fmt.Printf("First separator ranks: %v\n", S)
	// start := time.Now()

	// for isTwoArraysDiff(S, S_prev) {
	// 	S_prev = getValues(S)
	// 	for i := 0; i <= N-2; i++ {
	// 		start_v := max(0, S[max(1, i)-1])
	// 		end_v := min(S[i+1], N*C)
	// 		fmt.Printf("v from %v to %v\n", start_v, end_v)
	// 		for v := start_v; v <= end_v; v++ {
	// 			S_tmp = getValues(S)
	// 			S_tmp[i] = v
	// 			S_tmp[N-1] = calculateTail(S_tmp, N, C)
	// 			T_est := manager.estimate_traffic(S_tmp, numUsrs, numServers, numContents)
	// 			if T_est < T_min {
	// 				T_min = T_est
	// 				S = getValues(S_tmp)
	// 			}
	// 		}
	// 	}
	// }
	// t := time.Now()
	// elapsed := t.Sub(start)
	// fmt.Printf("%s ", "Elapsed time(s): ")
	// fmt.Println(elapsed)
	// fmt.Println(S)
	// return S

	// return make([]int, 0)
	return []int{36, 42, 54, 268}
}

func isTwoArraysDiff(a, b []int) bool {
	len := len(a)
	for i := 0; i < len; i++ {
		if a[i] != b[i] {
			return true
		}
	}
	return false
}

func getValues(a []int) []int {
	len := len(a)
	out := make([]int, len)
	for i := 0; i < len; i++ {
		out[i] = a[i]
	}
	return out
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func calculateTail(S []int, N, C int) int {
	result := C * N
	end_i := len(S) - 1
	for i := 0; i < end_i; i++ {
		if i == 0 {
			result -= S[i] * N
		} else {
			result -= (S[i] - S[i-1]) * (N - i)
		}
	}
	return result
}

func (manager *SpectrumManager_t) estimate_traffic(S_tmp []int, numUsrs, numServers, numContents int) float64 {
	return manager.network.expectTraffic().totalTraffic()
}

func (manager *SpectrumManager_t) BestReferenceRanks(mirageStore utils.MirageStore) []int {
	return make([]int, 0)
}

func isInArray(array []uint64, element uint64) bool {
	exist := false
	for _, v := range array {
		if element == v {
			exist = true
			break
		}
	}
	return exist
}

func (manager *SpectrumManager_t) adjacentSpectrums(node *Node_t) []uint64 {
	adjacentSpectrums := make([]uint64, 0)
	for _, link := range node.inputLinks {
		color, ok := manager.serverSpectrums[link.src]
		if ok {
			if !isInArray(adjacentSpectrums, color) {
				adjacentSpectrums = append(adjacentSpectrums, color)
			}
		}
	}
	for _, link := range node.outputLinks {
		color, ok := manager.serverSpectrums[link.dst]
		if ok {
			if !isInArray(adjacentSpectrums, color) {
				adjacentSpectrums = append(adjacentSpectrums, color)
			}
		}
	}
	return adjacentSpectrums
}

func (manager *SpectrumManager_t) availableSpectrums(node *Node_t) []uint64 {
	availableSpectrums := make([]uint64, 0)
	adjacentSpectrums := manager.adjacentSpectrums(node)
	for _, v := range manager.baseSpectrums {
		if !isInArray(adjacentSpectrums, v.(uint64)) {
			availableSpectrums = append(availableSpectrums, v.(uint64))
		}
	}
	return availableSpectrums
}

func (manager *SpectrumManager_t) countSpectrums(spectrum []uint64) int {
	return len(spectrum)
}

func (manager *SpectrumManager_t) selectDistantSpectrum(srcNode *Node_t, availableSpectrums []uint64) uint64 {
	// array stores minimal distance to node with available colors
	distances := make([]int, len(availableSpectrums))
	for i := range distances {
		distances[i] = -1 // distance value not found
	}
	// list to do BFS
	list := list.New()
	hopCount := 1

	for i := 0; i < len(srcNode.outputLinks); i++ {
		list.PushBack(srcNode.outputLinks[i].dst)
		list.PushBack(hopCount)
	}

	for !isAllElePositive(distances) {
		isInAvaiSpectrums := true
		node := list.Remove(list.Front()).(*Node_t)
		distance := list.Remove(list.Front()).(int)
		nodeColor, exist := manager.serverSpectrums[node]
		if exist {
			for i := 0; i < len(availableSpectrums); i++ {
				if nodeColor == availableSpectrums[i] {
					distances[i] = distance
					break
				}
			}
			isInAvaiSpectrums = false
		}
		if !exist || !isInAvaiSpectrums {
			hopCount++
			for i := 0; i < len(node.outputLinks); i++ {
				list.PushBack(node.outputLinks[i].dst)
				list.PushBack(hopCount)
			}
		}
	}
	minDistance := distances[0]
	selectedColor := availableSpectrums[0]
	for i := 0; i < len(availableSpectrums); i++ {
		if distances[i] < minDistance {
			minDistance = distances[i]
			selectedColor = availableSpectrums[i]
		}
	}
	return selectedColor
}

func (manager *SpectrumManager_t) setServerSpectrums() {
	verticesDegrees := make([]vertexDegree, 0)
	for _, node := range manager.network.nodes {
		verticesDegrees = append(verticesDegrees, vertexDegree{node, uint64(len(node.inputLinks) + len(node.outputLinks))})
	}

	sort.Slice(verticesDegrees, func(i, j int) bool {
		return verticesDegrees[i].degree > verticesDegrees[j].degree
	})

	for i := range verticesDegrees {
		availableSpectrums := manager.availableSpectrums(verticesDegrees[i].node)
		// missing sort descendingly based on minimal distance
		selectedColor := uint64(manager.bitSize)
		if i >= manager.bitSize {
			selectedColor = manager.selectDistantSpectrum(verticesDegrees[i].node, availableSpectrums)
		} else {
			selectedColor = manager.baseSpectrums[i].(uint64)
		}
		manager.serverSpectrums[verticesDegrees[i].node] = selectedColor
	}

	// fmt.Println("\nServers' color")
	// for _, node := range manager.network.nodes {
	// 	fmt.Printf("%s %d\n", node.id, manager.serverSpectrums[node])
	// }
}

func (manager *SpectrumManager_t) inspectServerSpectrums() {
	fmt.Println(manager.serverSpectrums)
}

func (manager *SpectrumManager_t) SetContentSpectrums(separatorRanks []int) {
	for i := 0; i < len(separatorRanks); i++ {
		manager.contentSpectrums = append(manager.contentSpectrums, uint64(separatorRanks[i]))
	}

	// Initialize spectrum tag
	manager.spectrumTags = make(map[int][]uint64)
	manager.initSpectrumTags()

	fmt.Println("XXXXXXXX-spectrumTags-XXXXX")
	for i := 1; i < len(manager.spectrumTags); i *= 2 {
		fmt.Printf("%d ", i)
		fmt.Println(manager.spectrumTags[i])
	}
	// share spectrum tag with cache.go
	for _, client := range manager.network.clients {
		client.Upstream().Storage().(*CacheStorage).spectrumTags = manager.spectrumTags
	}
}

func (manager *SpectrumManager_t) initSpectrumTags() {
	numColors := manager.bitSize
	numTags := int(math.Exp2(float64(numColors)))

	// colorTags table
	// colorTags[3]: get tags with 3 bits colors on
	colorTags := make([][]int, numColors+1)
	for i := range colorTags {
		colorTags[i] = make([]int, 0)
	}

	for i := 0; i < numTags; i++ {
		if i != 0 {
			colorTags[numBitOnes(i, numColors)] = append(colorTags[numBitOnes(i, numColors)], i)
		}
	}

	for i := 0; i < numColors; i++ {
		start := 1
		if i != 0 {
			start = int(manager.contentSpectrums[i-1] + 1)
		}

		end := int(manager.contentSpectrums[i])
		for key, j := start, 0; key <= end; key++ {
			manager.spectrumTags[key] = convertNumToArrayNum(colorTags[numColors-i][j], numColors)
			j = (j + i) % len(colorTags[numColors-i])
		}
	}
}

// get the number of bit 1s of a from 1st position to bitRange's nd
func numBitOnes(num, bitRange int) int {
	mask := 1
	numBitOnes := 0
	for i := 0; i < bitRange; i++ {
		if num&mask > 0 {
			numBitOnes++
		}
		mask <<= 1
	}
	return numBitOnes
}

func convertNumToArrayNum(num, bitRange int) []uint64 {
	mask := 1
	array := make([]uint64, bitRange)
	for i := 0; i < bitRange; i++ {
		if num&mask > 0 {
			array[bitRange-i-1] = 1
		}
		mask <<= 1
	}
	return array
}

type vertexDegree struct {
	node   *Node_t
	degree uint64
}

func isAllElePositive(array []int) bool {
	for _, v := range array {
		if v <= 0 {
			return false
		}
	}
	return true
}
