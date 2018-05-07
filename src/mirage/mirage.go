package main

import (
	"fmt"
	graphLoader "graph/loader"
	"math"
	"math/rand"
	"os"
	"parser"
	"time"
	"utils"
)

func init() {
	parser.ParseArgs()
	rand.Seed(int64(time.Now().Nanosecond()))
}

func main() {
	if parser.Options.GraphFilename == "" {
		utils.DebugPrint(fmt.Sprintln("Please specify graph filename."))
		parser.PrintDefaults()
		os.Exit(1)
	}

	network, _ := graphLoader.LoadGraph(parser.Options.GraphFilename)
	warmupRequestCount := parser.Options.WarmupRequestCount
	evaluationRequestCount := parser.Options.EvaluationRequestCount

	if network.CacheAlgorithm() == "iris" {
		bitSize := parser.Options.SpectrumBitSize
		librarySize := network.LibrarySize()
		spectrumCapacity := network.SpectrumCapacity()

		var bestSeparatorRanks []int
		mirageStore := utils.LoadData("mirage.dat")

		fmt.Println("XXXXX-mirageStore")
		fmt.Println(mirageStore)

		if parser.Options.UseReferenceRanks {
			utils.DebugPrint(fmt.Sprintln("deciding separator ranks from references..."))
			bestSeparatorRanks = network.SpectrumManager().BestReferenceRanks(mirageStore)
			utils.DebugPrint(fmt.Sprintln(bestSeparatorRanks))
		} else {
			for _, best := range mirageStore.BestSeparatorRanks {
				if network.MatchStoreData(best, bitSize) {
					bestSeparatorRanks = best.Ranks
					break
				}
			}

			if bestSeparatorRanks == nil {
				bestSeparatorRanks = make([]int, 4)
				bestSeparatorRanks[len(bestSeparatorRanks)-1] = int(math.Min(float64(spectrumCapacity*4), float64(librarySize)))
				bestSeparatorRanks = network.SpectrumManager().BestSeparatorRanks(bestSeparatorRanks)

				fmt.Println("XXXXX-bestSeperatorRanks")
				fmt.Println(bestSeparatorRanks)

				best := network.GenerateBestSeparatorRanksData(bestSeparatorRanks)
				mirageStore.BestSeparatorRanks = append(mirageStore.BestSeparatorRanks, best)
				utils.StoreData("mirage.dat", mirageStore)
			}
		}

		for len(bestSeparatorRanks) < bitSize {
			separatorRanks := make([]int, len(bestSeparatorRanks)*2)
			for index := range bestSeparatorRanks {
				separatorRanks[index*2] = bestSeparatorRanks[index]
			}
			bestSeparatorRanks = network.SpectrumManager().BestSeparatorRanks(separatorRanks)
		}

		network.SpectrumManager().SetContentSpectrums(bestSeparatorRanks)
	}

	/*
		// CUSTOM CODE

		// make a map for easy access to the client variables by their upstream node IDs.
		clients := make(map[string]graph.Client)
		for _, client := range network.Clients() {
			clients[client.Upstream().ID()] = client
		}

		// read log file
		fileHandle, _ := os.Open("data.txt")
		defer fileHandle.Close()
		fileScanner := bufio.NewScanner(fileHandle)

		var data, sumData []int
		totalID, totalCount := 0, 0

		for fileScanner.Scan() {
			s := strings.Split(fileScanner.Text(), ",")
			count, _ := strconv.Atoi(s[1])
			totalCount += count
			data = append(data, count)
			sumData = append(sumData, totalCount)
			totalID++
		}

		fmt.Printf("Total ID: %d\nTotal count: %d\n", totalID, totalCount)
		// fmt.Println(data)
		// fmt.Println(sumData)
		// fmt.Println(clients)
		// fmt.Println(network.Clients())

		// generate requests for warming up
		for index := 0; index < warmupRequestCount; index++ {
			utils.DebugPrint(fmt.Sprintf("\rwarming up: (%d/%d)", index+1, warmupRequestCount))
			for node := 1; node <= len(clients); node++ {
				randNum := rand.Intn(totalCount-1) + 1
				for id := 0; id < totalID; id++ {
					if randNum > sumData[id] {
						continue
					}
					clients[fmt.Sprintf("node%d", node)].RequestByID(id)
					break
				}
			}
		}

		utils.DebugPrint(fmt.Sprintln())
		network.ResetCounters()

		// generate requests
		for index := 0; index < evaluationRequestCount; index++ {
			utils.DebugPrint(fmt.Sprintf("\revaluation: (%d/%d)", index+1, evaluationRequestCount))
			for node := 1; node <= len(clients); node++ {
				randNum := rand.Intn(totalCount-1) + 1
				for id := 0; id < totalID; id++ {
					if randNum > sumData[id] {
						continue
					}
					clients[fmt.Sprintf("node%d", node)].RequestByID(id)
					break
				}
			}
		}

		for index := 0; index < evaluationRequestCount; index++ {
			utils.DebugPrint(fmt.Sprintf("\revaluation: (%d/%d)", index+1, evaluationRequestCount))
			for node := 1; node <= len(clients); node++ {
				randNum := rand.Intn(totalCount-1) + 1
				for id := 0; id < totalID; id++ {
					if randNum > sumData[id] {
						continue
					}
					clients[fmt.Sprintf("node%d", node)].RequestByID(id)
					break
				}
			}
		}

		// END CUSTOM CODE
	*/

	// generate requests for warming up
	for index := 0; index < warmupRequestCount; index++ {
		utils.DebugPrint(fmt.Sprintf("\rwarming up: (%d/%d)", index+1, warmupRequestCount))
		for _, client := range network.Clients() {
			client.RandomRequest()
		}
	}

	utils.DebugPrint(fmt.Sprintln())
	network.ResetCounters()

	// generate requests
	for index := 0; index < evaluationRequestCount; index++ {
		utils.DebugPrint(fmt.Sprintf("\revaluation: (%d/%d)", index+1, evaluationRequestCount))
		for _, client := range network.Clients() {
			client.RandomRequest()
		}
	}

	utils.DebugPrint(fmt.Sprintln())

	switch parser.Options.OutputFormat {
	case "plain":
		totalTraffic := 0.0
		for _, link := range network.Links() {
			totalTraffic += link.Traffic()
		}
		originTraffic := 0.0
		for _, link := range network.OriginServers()[0].OutputLinks() {
			originTraffic += link.Traffic()
		}
		internalTraffic := totalTraffic - originTraffic

		utils.DebugPrint(fmt.Sprintln("--"))
		utils.DebugPrint(fmt.Sprintln("request_count: ", evaluationRequestCount*len(network.Clients())))
		utils.DebugPrint(fmt.Sprintln("total_traffic: ", totalTraffic))
		utils.DebugPrint(fmt.Sprintln(" - internal_traffic:", internalTraffic))
		utils.DebugPrint(fmt.Sprintln(" - origin_traffic:  ", originTraffic))
		utils.DebugPrint(fmt.Sprintln("--"))

		network.PlainReport()
	case "json":
		network.JsonReport()
	}
}
