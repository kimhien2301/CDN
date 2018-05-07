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
			// load from mirage.dat
			for _, best := range mirageStore.BestSeparatorRanks {
				if network.MatchStoreData(best, bitSize) {
					bestSeparatorRanks = best.Ranks
					break
				}
			}

			// calculate separator ranks and store to mirage.dat
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
