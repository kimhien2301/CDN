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
	baseSpectrums        []interface{}                         // colors not sure ???
	spectrumTags         map[int][]uint64                      // contents' color
	bitSize              int                                   // # colors
	contentSpectrums     []uint64                              // seperatorRanks
	serverSpectrums      map[*Node_t]uint64                    // servers' color
	spectrumRoutingTable map[*Node_t]map[uint64][]ForwardEntry // color-based routing table
}

// serverSpectrums      map[*Node_t]uint64
// Server color 	Color Tag
// 0				0001
// 1				0010
// 2				0100
// 3				1000

func newSpectrumManager(bitSize int, network *Graph_t) *SpectrumManager_t {
	manager := new(SpectrumManager_t)
	// TU ADD
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
	// TU add
	manager.spectrumRoutingTable = make(map[*Node_t]map[uint64][]ForwardEntry)
	serverSpectrums := manager.serverSpectrums
	routingTable := manager.network.router.RoutingTable()

	for _, node := range manager.network.nodes {
		// Create Request Routing Table for a node
		nodeRequestRoutingTable := make(map[uint64][]ForwardEntry)
		queue := list.New()
		enqueuedEles := list.New()

		for _, adjacentNode := range node.outputAdjacent() {
			queue.PushBack(adjacentNode)
			enqueuedEles.PushBack(adjacentNode)
		}

		// fmt.Printf("%s %s\n", "Traversed node", node.ID())
		// fmt.Println("Inspect queue")
		// for ele := queue.Front(); ele != nil; ele = ele.Next() {
		// 	fmt.Println(ele.Value.(*Node_t).ID())
		// }

		for queue.Len() != 0 && len(nodeRequestRoutingTable) < 4 {
			ele := queue.Front()
			adjNode := queue.Remove(ele).(*Node_t)

			_, exist := nodeRequestRoutingTable[serverSpectrums[adjNode]]

			if !exist && adjNode != node {
				forwardEntries := routingTable[node.id][adjNode.id]
				var forwardEntryWithMinCost ForwardEntry

				for i := range forwardEntries {
					if forwardEntryWithMinCost == nil {
						forwardEntryWithMinCost = forwardEntries[i]
					} else {
						if forwardEntryWithMinCost.Cost() > forwardEntries[i].Cost() {
							forwardEntryWithMinCost = forwardEntries[i]
						}
					}
				}

				nodeRequestRoutingTable[serverSpectrums[adjNode]] = append(nodeRequestRoutingTable[serverSpectrums[adjNode]], forwardEntryWithMinCost)
			}

			// Insert new adjacent nodes of current checking node
			for _, adjacentNode := range adjNode.outputAdjacent() {
				newAdjNodeExist := false
				for ele := enqueuedEles.Front(); ele != nil; ele = ele.Next() {
					if ele.Value.(*Node_t) == adjacentNode {
						newAdjNodeExist = true
						break
					}
				}

				if !newAdjNodeExist {
					queue.PushBack(adjacentNode)
					enqueuedEles.PushBack(adjacentNode)
				}
			}

			// fmt.Println("Inspect queue")
			// for ele := queue.Front(); ele != nil; ele = ele.Next() {
			// 	fmt.Println(ele.Value.(*Node_t).ID())
			// }
		}
		manager.spectrumRoutingTable[node] = nodeRequestRoutingTable
	}

}

func (manager *SpectrumManager_t) inspectSpectrumRoutingTable() {
	spectrumRoutingTable := manager.spectrumRoutingTable
	fmt.Println("\nRequest Routing Table")
	for _, node := range manager.network.nodes {
		fmt.Println(spectrumRoutingTable[node])
	}
	fmt.Printf("%8s %10s %s\n", "Node ID", "Color tag", "ForwardEntry")
	for _, node := range manager.network.nodes {
		for i := 0; i < 4; i++ {
			if len(spectrumRoutingTable[node][uint64(i)]) > 0 {
				fmt.Printf("%8s %10d ", node.id, i)
				fmt.Printf("%6s %f\n", spectrumRoutingTable[node][uint64(i)][0].Node().ID(),
					spectrumRoutingTable[node][uint64(i)][0].Cost())
			}
		}
		fmt.Println(spectrumRoutingTable[node])
	}
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

// TU modified code
func (manager *SpectrumManager_t) BestSeparatorRanks(separatorRanks []int) []int {
	// N := manager.bitSize                                            // # colors
	// C := manager.network.clients[0].Upstream().Storage().Capacity() // cache server capacity
	// numUsrs := len(manager.network.clients)
	// numServers := len(manager.network.nodes)
	// numContents := manager.network.LibrarySize()

	// fmt.Println("\nnumUsrs numServers numContents")
	// fmt.Println(numUsrs)
	// fmt.Println(numServers)
	// fmt.Println(numContents)

	// S := make([]int, N)
	// S_prev := make([]int, N)
	// var S_tmp []int
	// T_min := math.MaxFloat64

	// S[N-1] = N * C
	// fmt.Println("Calculate separator ranks")
	// fmt.Println(S)
	// start := time.Now()

	// for isTwoArraysDiff(S, S_prev) {
	// 	S_prev = getValues(S)
	// 	for i := 0; i <= N-2; i++ {
	// 		start_v := max(0, S[max(1, i)-1])
	// 		end_v := min(S[i+1], N*C)
	// 		for v := start_v; v <= end_v; v++ {
	// 			S_tmp = getValues(S)
	// 			S_tmp[i] = v
	// 			S_tmp[N-1] = calculateTail(S_tmp, N, C)
	// 			T_est := manager.estimate_traffic(S_tmp, numUsrs, numServers, numContents)
	// 			if T_est < T_min {
	// 				T_min = T_est
	// 				S = getValues(S_tmp)
	// 				fmt.Println(S)
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

	// traffic := 0.0
	// routingTable := est manager.network.router.(*Router_t).routingTable
	// graph := manager.network

	// dist := graph.clients[0].Dist()

	// for i := 0; i < numUsrs; i++ {
	// 	for j := 0; j < numServers; j++ {
	// 		src := graph.clients[i].Upstream().ID()
	// 		des := graph.nodes[j].id
	// 		tmp := routingTable[src][des]
	// 		cost := 0.0
	// 		if len(tmp) == 0 {
	// 			cost = 1.0
	// 		} else {
	// 			cost = tmp[0].Cost() + 1
	// 		}
	// 		for k := 0; k < numContents; k++ {
	// 			bin_var := 0.0
	// 			if graph.clients[i].Upstream().Storage().Exist(k + 1) {
	// 				bin_var = 1.0
	// 			} else {
	// 				graph.clients[i].Upstream().Storage().Insert(k+1, k+1)
	// 			}

	// 			traffic += cost * dist.PDF(k+1) * bin_var
	// 			//fmt.Println(cost * dist.PDF(k+1) * bin_var)
	// 		}
	// 	}
	// }
	// return traffic
	// fmt.Println(manager.network.expectTraffic().totalTraffic())
	return manager.network.expectTraffic().totalTraffic()
}

//
func (manager *SpectrumManager_t) BestReferenceRanks(mirageStore utils.MirageStore) []int {
	return make([]int, 0)
}

// TU modified code
func isInArray(array []uint64, elem uint64) bool {
	exist := false
	for _, v := range array {
		if elem == v {
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
	hopcount := 1

	for i := 0; i < len(srcNode.outputLinks); i++ {
		list.PushBack(srcNode.outputLinks[i].dst)
		list.PushBack(hopcount)
	}

	// fmt.Println(ele.Value)
	for !isAllElePositive(distances) {
		isInAvaiSpectrums := true

		node := list.Remove(list.Front()).(*Node_t)
		//fmt.Println(node.id)
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
			hopcount++
			for i := 0; i < len(node.outputLinks); i++ {
				list.PushBack(node.outputLinks[i].dst)
				list.PushBack(hopcount)
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
		verticesDegrees = append(verticesDegrees,
			vertexDegree{node, uint64(len(node.inputLinks) + len(node.outputLinks))})
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

	// fmt.Println("XXXXXXXXXX-spectrumTags-XXXXX")
	// for i := 1; i < len(manager.spectrumTags); i *= 2 {
	// 	fmt.Printf("%d ", i)
	// 	fmt.Println(manager.spectrumTags[i])
	// }
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
			j = (j + 1) % len(colorTags[numColors-i])
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

func convertNumToArrayNum(num int, bitRange int) []uint64 {
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

func (manager *SpectrumManager_t) SetContentWithTag(contentID int, tag []uint64) {
	manager.spectrumTags[contentID] = tag
	// fmt.Printf("Content %d tag ", contentID)
	// fmt.Println(tag)
}

// TU add code
type vertexDegree struct {
	node   *Node_t
	degree uint64
}

// type listElement {
// 	node *Node_t
// 	distance int
// }

func isAllElePositive(array []int) bool {
	for _, v := range array {
		if v <= 0 {
			return false
		}
	}
	return true
}
