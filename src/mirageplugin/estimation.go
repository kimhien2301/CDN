package main

import "fmt"
import "graph"
import "parser"
import "cache"
import "utils"
import "algorithm"

type TrafficExpectation_t struct {
	network *Graph_t
	Traffic map[graph.UnidirectionalLink]float64
}

func newTrafficExpectation(network *Graph_t) *TrafficExpectation_t {
	expectation := new(TrafficExpectation_t)
	expectation.network = network
	expectation.Traffic = make(map[graph.UnidirectionalLink]float64)
	for _, link := range network.links {
		expectation.Traffic[link] = 0.0
	}
	return expectation
}

func (expectation *TrafficExpectation_t) totalTraffic() float64 {
	totalTraffic := 0.0
	for _, link := range expectation.network.links {
		totalTraffic += expectation.Traffic[link]
	}
	return totalTraffic
}

func (expectation *TrafficExpectation_t) originTraffic() float64 {
	originTraffic := 0.0
	for _, origin := range expectation.network.originServers() {
		for _, link := range origin.outputLinks {
			originTraffic += expectation.Traffic[link]
		}
	}
	return originTraffic
}

func (expectation *TrafficExpectation_t) internalTraffic() float64 {
	return expectation.totalTraffic() - expectation.originTraffic()
}

func (expectation *TrafficExpectation_t) inspect() {
	utils.DebugPrint(fmt.Sprintln("--"))
	for _, link := range expectation.network.links {
		utils.DebugPrint(fmt.Sprintf("%s -> %s :\t%f\n", link.src.ID, link.dst.ID, expectation.Traffic[link]))
	}
	utils.DebugPrint(fmt.Sprintln("--"))
	linkTraffic := make([]float64, 0)
	for _, traffic := range expectation.Traffic {
		linkTraffic = append(linkTraffic, traffic)
	}
	utils.DebugPrint(fmt.Sprintf("Link traffic SD: %f\n", algorithm.StandardDeviation(linkTraffic)))
	utils.DebugPrint(fmt.Sprintln("--"))
	utils.DebugPrint(fmt.Sprintln("total_traffic: ", expectation.totalTraffic()))
	utils.DebugPrint(fmt.Sprintln(" - internal_traffic:", expectation.internalTraffic()))
	utils.DebugPrint(fmt.Sprintln(" - origin_traffic:  ", expectation.originTraffic()))
}

func (network *Graph_t) expectTraffic() *TrafficExpectation_t {
	var forwardEntry ForwardEntry
	network.fillCaches()
	expectation := newTrafficExpectation(network)

	// Generate requests for all contents from each client
	for _, client := range network.clients {
		for contentRank := 1; contentRank <= network.originServers()[0].Entity().(*ServerModel_t).Storage().Len(); contentRank++ {
			fromID := client.Upstream().ID()
			// Set the amount of traffic generated in one transfer
			unitTraffic := client.Dist().PDF(contentRank)
			// Generate request
			request := cache.ContentRequest{contentRank, make([]interface{}, 0), client.TrafficWeight(), make([]string, 0)}
			// If the server hits on the client immediately, the traffic volume = 0
			if client.Upstream().Storage().Exist(request.ContentKey) {
				continue
			}

			for {
				// Update server that passed through
				loopDetected := false
				for _, nodeID := range request.XForwardedFor {
					if nodeID == fromID {
						loopDetected = true
					}
				}
				request.XForwardedFor = append(request.XForwardedFor, fromID)

				// Determine the next transfer destination (change the route control method by the cache algorithm）
				if network.model.Nodes[0].CacheAlgorithm == "iris" && loopDetected == false && !parser.Options.UseShortestPath {
					forwardEntry = network.router.selectForwardEntryBySpectrum(fromID, request)
				} else {
					forwardEntry = network.router.selectForwardEntry(fromID, request)
				}

				// Add traffic to the transfer link
				expectation.Traffic[forwardEntry.Link()] += unitTraffic

				// When the next node hits, routing control end
				if forwardEntry.Node().Entity().(*ServerModel_t).Storage().Exist(request.ContentKey) {
					break
				}

				// Make the request originator the next server and forward the request again
				fromID = forwardEntry.Node().ID()
			}
		}
	}

	return expectation
}

func (network *Graph_t) fillCaches() {
	switch network.model.Nodes[0].CacheAlgorithm {
	case "iris":
		// 	for _, node := range network.cacheServers() {
		// 		node.Entity().(*ServerModel_t).Storage().(iris.Accessor).FillUp()
		// 	}
		for i := 0; i < 5000; i++ {
			for _, client := range network.Clients() {
				client.RandomRequest()
			}
		}
	}
}

// ADD
func (network *Graph_t) expectTrafficWithNoFillCaches() *TrafficExpectation_t {
	var forwardEntry ForwardEntry
	//network.fillCaches()
	expectation := newTrafficExpectation(network)

	// 各クライアントからすべてのコンテンツについてのリクエストを生成
	for _, client := range network.clients {
		for contentRank := 1; contentRank <= network.originServers()[0].Entity().(*ServerModel_t).Storage().Len(); contentRank++ {
			fromID := client.Upstream().ID()
			// 1度の転送に発生する通信量を設定
			unitTraffic := client.Dist().PDF(contentRank)
			// リクエストを生成
			request := cache.ContentRequest{contentRank, make([]interface{}, 0), client.TrafficWeight(), make([]string, 0)}
			// クライアント直上のサーバでヒットしたら通信量=0
			if client.Upstream().Storage().Exist(request.ContentKey) {
				continue
			}

			for {
				// 経由したサーバを更新
				loopDetected := false
				for _, nodeID := range request.XForwardedFor {
					if nodeID == fromID {
						loopDetected = true
					}
				}
				request.XForwardedFor = append(request.XForwardedFor, fromID)

				// 次の転送先を決定（キャッシュアルゴリズムによって経路制御方法を変える）
				if network.model.Nodes[0].CacheAlgorithm == "iris" && loopDetected == false && !parser.Options.UseShortestPath {
					forwardEntry = network.router.selectForwardEntryBySpectrum(fromID, request)
				} else {
					forwardEntry = network.router.selectForwardEntry(fromID, request)
				}

				// 転送リンクに通信量を加算
				expectation.Traffic[forwardEntry.Link()] += unitTraffic

				// 次のノードでヒットしたら経路制御終了
				if forwardEntry.Node().Entity().(*ServerModel_t).Storage().Exist(request.ContentKey) {
					break
				}

				// リクエスト発生元を次のサーバにしてもう一度リクエストを転送
				fromID = forwardEntry.Node().ID()
			}
		}
	}

	return expectation
}
