package main

import (
	"bufio"
	"fmt"
	"graph"
	graphLoader "graph/loader"
	"math"
	"math/rand"
	"os"
	"parser"
	"sort"
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
		// fmt.Printf("Spectrum capacity: %v\n", spectrumCapacity)

		var bestSeparatorRanks []int
		mirageStore := utils.LoadData("mirage.dat")
		// fmt.Printf("Mirage Store: %v\n", mirageStore)

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

				fmt.Printf("Best Seperator Ranks: %v\n", bestSeparatorRanks)

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

	// fmt.Printf("Cache Algorithm: %v\n", network.CacheAlgorithm())
	if !parser.Options.UseShortestPath {
		fmt.Println("Color-based routing")
	}
	//---------- GA ----------//
	// parser.GA = parser.Options.GA
	if parser.Options.GA {
		fmt.Println("Genetic Algorithm")

		totalContent := network.LibrarySize()
		cacheCapacity := network.GetCacheCapacity()
		numCacheServer := network.GetNumberOfCacheServers()
		numChromosome := 100
		generation := 5
		mutanRate := float64(0.01)
		parent := initializePopulation(network, totalContent, cacheCapacity, numChromosome, numCacheServer)
		offspring := initializePopulation(network, totalContent, cacheCapacity, numChromosome, numCacheServer)
		parentFitness := computeFitness(network, parent)
		// offspringFitness := computeFitness(network, offspring)
		var offspringFitness []float64
		// crossover(offspring, parent)
		offspringFitness = computeFitness(network, offspring)

		var temp float64

		start := time.Now()
		for i := 0; i < generation; i++ {
			// crossover(offspring, parent)
			// offspringFitness = computeFitness(network, offspring)
			// selection(offspring, parent, offspringFitness, parentFitness)
			// mutation(offspring, mutanRate, cacheCapacity, totalContent)
			// offspringFitness = computeFitness(network, offspring)
			// selection(offspring, parent, offspringFitness, parentFitness)

			selection(offspring, parent, offspringFitness, parentFitness)
			crossover(offspring, parent)
			mutation(network, offspring, mutanRate, cacheCapacity, totalContent)
			offspringFitness = computeFitness(network, offspring)

			if i == 0 || temp != offspringFitness[0] {
				fmt.Printf("Generation %d: %f\n", i, offspringFitness[0])
				//utils.DebugPrint(fmt.Sprintf("\rGeneration %d: %f", i, offspringFitness[0]))
			} else {
				utils.DebugPrint(fmt.Sprintf("\rGeneration %d: %f - ", i, offspringFitness[0]))
			}

			temp = offspringFitness[0]
		}
		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Printf("%s ", "Generations: ")
		fmt.Println(generation)
		fmt.Printf("%s ", "Fitness: ")
		fmt.Println(offspringFitness[0])
		fmt.Printf("%s ", "Elapsed time GA(s): ")
		fmt.Println(elapsed)

		// GA log to file
		f, err := os.Create("GA_log_file.txt")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		w := bufio.NewWriter(f)
		defer w.Flush()
		fmt.Fprintf(w, "Generations: ")
		fmt.Fprintln(w, generation)
		fmt.Fprintf(w, "Elapsed time GA(s): ")
		fmt.Fprintln(w, elapsed)
		fmt.Fprintf(w, "Fitness: ")
		fmt.Fprintln(w, offspringFitness[0])
		fmt.Fprintf(w, "Solution: ")
		fmt.Fprintln(w, offspring[0])

		network.SetCacheServers(offspring[0])

		// genServerListOfContents(network, cacheCapacity, totalContent)
	} else if parser.Options.InsertNewContents {
		fmt.Println("Insert New Contents")
		// fmt.Printf("%p \n", network)
		// network2, _ := graphLoader.LoadGraph(parser.Options.GraphFilename)
		// fmt.Printf("%p \n", network2)
		//fmt.Println(network2.OriginServers()[0].Entity().(graph.ServerModel).Storage().CacheList())

		network.SpectrumManager().SetContentWithTag(1, []uint64{0, 0, 0, 0})
		network.SpectrumManager().SetContentWithTag(2, []uint64{0, 0, 0, 0})
		network.SpectrumManager().SetContentWithTag(3, []uint64{0, 0, 0, 0})
		network.SpectrumManager().SetContentWithTag(4, []uint64{0, 0, 0, 0})
		network.SpectrumManager().SetContentWithTag(5, []uint64{0, 0, 0, 0})

		// network.ViewCacheContents()
		for index := 0; index < warmupRequestCount; index++ {
			utils.DebugPrint(fmt.Sprintf("\rwarming up: (%d/%d)", index+1, warmupRequestCount))
			for _, client := range network.Clients() {
				client.RandomRequestForInsertNewContents()
			}
		}
		// network.ViewCacheContents()

		utils.DebugPrint(fmt.Sprintln())
		// network.ResetCounters()

		fmt.Println("Evaluation Phase")

		f, err := os.Create("InsertNewContents.txt")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		w := bufio.NewWriter(f)
		defer w.Flush()

		fmt.Fprintf(w, "Evaluation number cache hit rate\n")
		fmt.Fprintf(w, "0,%5.3f\n", network.GetCacheHitRate())

		for index := 0; index < 500; index++ {
			utils.DebugPrint(fmt.Sprintf("\revaluation: (%d/%d)", index+1, evaluationRequestCount))
			network.ResetCounters()
			for _, client := range network.Clients() {
				client.RandomRequestForInsertNewContents()
			}
			fmt.Fprintf(w, "%d,%5.3f\n", index+1, network.GetCacheHitRate())
		}
		network.PlainReport()
		// network.ViewCacheContents()
		// A
		for index := 0; index < 500; index++ {
			utils.DebugPrint(fmt.Sprintf("\revaluation: (%d/%d)", index+501, evaluationRequestCount))
			network.ResetCounters()
			for _, client := range network.Clients() {
				client.RandomRequest()
			}
			fmt.Fprintf(w, "%d,%5.3f\n", index+501, network.GetCacheHitRate())
		}

		// network.ViewCacheContents()
		// B
		// Insert 5 new contents
		// for index := 0; index < 5; index++ {
		// 	utils.DebugPrint(fmt.Sprintf("\revaluation: (%d/%d)", index+501, evaluationRequestCount))

		// 	for _, client := range network.Clients() {
		// 		client.RequestByID(index + 1)
		// 	}

		// 	fmt.Fprintf(w, "%d,%5.3f\n", index+501, network.GetCacheHitRate())

		// }
		// network.ViewCacheContents()
		// for _, client := range network.Clients() {
		// 	client.RequestByID(0 + 1)
		// }
		// for _, client := range network.Clients() {
		// 	client.RequestByID(1 + 1)
		// }
		// network.PlainReport()
		// for index := 5; index < 500; index++ {
		// 	utils.DebugPrint(fmt.Sprintf("\revaluation: (%d/%d)", index+501, evaluationRequestCount))
		// 	// network.ResetCounters()
		// 	for _, client := range network.Clients() {
		// 		client.RandomRequest()
		// 	}
		// 	fmt.Fprintf(w, "%d,%5.3f\n", index+501, network.GetCacheHitRate())
		// }

		// network.PlainReport()

		// utils.DebugPrint(fmt.Sprintln())

		// totalTraffic := 0.0
		// for _, link := range network.Links() {
		// 	totalTraffic += link.Traffic()
		// }
		// originTraffic := 0.0
		// for _, link := range network.OriginServers()[0].OutputLinks() {
		// 	originTraffic += link.Traffic()
		// }
		// internalTraffic := totalTraffic - originTraffic

		// utils.DebugPrint(fmt.Sprintln("--"))
		// utils.DebugPrint(fmt.Sprintln("request_count: ", evaluationRequestCount*len(network.Clients())))
		// utils.DebugPrint(fmt.Sprintln("total_traffic: ", totalTraffic))
		// utils.DebugPrint(fmt.Sprintln(" - internal_traffic:", internalTraffic))
		// utils.DebugPrint(fmt.Sprintln(" - origin_traffic:  ", originTraffic))
		// utils.DebugPrint(fmt.Sprintln("--"))

	} else {
		// generate requests for warming up
		for index := 0; index < warmupRequestCount; index++ {
			utils.DebugPrint(fmt.Sprintf("\rwarming up: (%d/%d)", index+1, warmupRequestCount))
			for _, client := range network.Clients() {
				// fmt.Println(client.Upstream().ID())
				client.RandomRequest()
			}
		}

		utils.DebugPrint(fmt.Sprintln())
		network.ResetCounters()

		// generate requests
		for index := 0; index < evaluationRequestCount; index++ {
			utils.DebugPrint(fmt.Sprintf("\revaluation: (%d/%d)", index+1, evaluationRequestCount))
			for _, client := range network.Clients() {
				// fmt.Println(client.Upstream().ID())
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
			/*
				fmt.Println("--")
				fmt.Println("request_count: ", evaluationRequestCount*len(network.Clients()))
				fmt.Println("total_traffic: ", totalTraffic)
				fmt.Println(" - internal_traffic:", internalTraffic)
				fmt.Println(" - origin_traffic:  ", originTraffic)
				fmt.Println("--")
			*/
			network.PlainReport()
		case "json":
			network.JsonReport()

		}
	}
}

func logToFile() {
	f, err := os.Create("GA_log_file.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	for j := 0; j < 5; j++ {
		count := 4 //count := rand.Intn(100)
		for i := 0; i < count; i++ {
			fmt.Fprint(w, rand.Intn(1000), " ")
		}
		fmt.Fprintln(w)
	}
}

type Chromosome [][]int
type Population []Chromosome

func genServerListOfContents(network graph.Graph, cacheCapacity int) []int {
	cacheMap := make(map[int]int)

	for {
		key := network.Clients()[0].Dist().Intn()
		_, exist := cacheMap[key]
		if exist {
			cacheMap[key]++
		} else {
			cacheMap[key] = 1
		}

		if len(cacheMap) == cacheCapacity {
			break
		}
	}

	returnedArray := make([]int, 0)

	for i := range cacheMap {
		returnedArray = append(returnedArray, i)
	}

	// fmt.Printf("Cache map:\n %v\n", cacheMap)
	// fmt.Printf("Returned array:\n %v\n", returnedArray)

	return returnedArray
}

func initializeChromosome(network graph.Graph, individual Chromosome, cacheCapacity, numCacheServer, totalContent int) {
	for i := 0; i < numCacheServer; i++ {
		individual[i] = make([]int, cacheCapacity)

		values := genServerListOfContents(network, cacheCapacity)
		for j := 0; j < cacheCapacity; j++ {
			individual[i][j] = values[j] + 1
		}
	}
}

func initializePopulation(network graph.Graph, totalContent, cacheCapacity, numChromosome, numCacheServer int) Population {
	population := make([]Chromosome, numChromosome)
	for i := 0; i < numChromosome; i++ {
		population[i] = make(Chromosome, numCacheServer)
		initializeChromosome(network, population[i], cacheCapacity, numCacheServer, totalContent)
	}
	return population
}

func computeFitness(network graph.Graph, population Population) []float64 {
	lenPop := len(population)
	fitness := make([]float64, lenPop)

	for i := 0; i < lenPop; i++ {
		network.SetCacheServers(population[i])
		fitness[i] = network.GetExpectTraffic() // fix fill cache in getExpectTraffic()
	}
	return fitness
}

func crossover(offspring, parent Population) {
	sizePop := len(offspring)
	numCacheServers := len(offspring[0])  // number of cache servers
	cacheCapacity := len(offspring[0][0]) // cache capacity

	for i := 0; i < sizePop; i++ {
		crossPoint := rand.Intn(numCacheServers)
		crossChromosome := rand.Intn(sizePop)

		for j := crossPoint; j < numCacheServers; j++ {
			for k := 0; k < cacheCapacity; k++ {
				offspring[i][j][k] = parent[crossChromosome][j][k]
			}
		}
	}
}

func mutation(network graph.Graph, population Population, mutanRate float64, cacheCapacity, totalContent int) {
	sizePop := len(population)
	numCacheServers := len(population[0]) // number of cache servers

	for i := 0; i < sizePop; i++ {
		for j := 0; j < numCacheServers; j++ {
			if rand.Float64() < mutanRate {
				population[i][j] = make([]int, cacheCapacity)

				values := genServerListOfContents(network, cacheCapacity)
				for k := 0; k < cacheCapacity; k++ {
					population[i][j][k] = values[k] + 1
				}
			}
		}
	}
}

// Sort value by key in increasing order
type SortByKey struct {
	value Population
	key   []float64
}

func (s SortByKey) Len() int {
	return len(s.value)
}

func (s SortByKey) Swap(x, y int) {
	//sizePop := len(s.value)
	numCacheServers := len(s.value[0])  // number of cache servers
	cacheCapacity := len(s.value[0][0]) // cache capacity

	for j := 0; j < numCacheServers; j++ {
		for k := 0; k < cacheCapacity; k++ {
			s.value[x][j][k], s.value[y][j][k] = s.value[y][j][k], s.value[x][j][k]
		}
	}

	s.key[x], s.key[y] = s.key[y], s.key[x]
}

func (s SortByKey) Less(i, j int) bool {
	return s.key[i] < s.key[j]
}

func cloneChromosome(chromosome Chromosome) Chromosome {
	numCacheServers := len(chromosome)
	cacheCapacity := len(chromosome[0])

	newChromosome := make(Chromosome, numCacheServers)
	for i := 0; i < numCacheServers; i++ {
		newChromosome[i] = make([]int, cacheCapacity)

		for j := 0; j < cacheCapacity; j++ {
			newChromosome[i][j] = chromosome[i][j]
		}
	}
	return newChromosome
}

func selection(offspring, parent Population, offspringFitness, parentFitness []float64) {
	offspringSorting := new(SortByKey)
	parentSorting := new(SortByKey)

	offspringSorting.key = offspringFitness
	offspringSorting.value = offspring
	parentSorting.key = parentFitness
	parentSorting.value = parent

	sort.Sort(offspringSorting)
	sort.Sort(parentSorting)

	sizePop := len(offspring)

	newPopulation := make(Population, sizePop)
	newFitness := make([]float64, sizePop)

	offspringPos := 0
	parentPos := 0

	for i := 0; i < sizePop; i++ {
		if offspringPos < sizePop && parentPos < sizePop {
			if offspringFitness[offspringPos] < parentFitness[parentPos] {
				newPopulation[i] = cloneChromosome(offspring[offspringPos])
				newFitness[i] = offspringFitness[offspringPos]
				offspringPos++
			} else {
				newPopulation[i] = cloneChromosome(parent[parentPos])
				newFitness[i] = parentFitness[parentPos]
				parentPos++
			}
		} else if offspringPos < sizePop {
			newPopulation[i] = cloneChromosome(offspring[offspringPos])
			newFitness[i] = offspringFitness[offspringPos]
			offspringPos++
		} else {
			newPopulation[i] = cloneChromosome(parent[parentPos])
			newFitness[i] = parentFitness[parentPos]
			parentPos++
		}
	}

	// fmt.Println("NEW POP")
	// fmt.Println(newPopulation)
	// fmt.Println(newFitness)

	//parent = newPopulation
	//parentFitness = newFitness

	//offspring = make(Population, sizePop)
	//offspringFitness = make([]float64, sizePop)

	for i, v := range newFitness {
		offspring[i] = cloneChromosome(newPopulation[i])
		offspringFitness[i] = v
		parent[i] = cloneChromosome(newPopulation[i])
		parentFitness[i] = v
	}
}

func gaTest(network graph.Graph) {
	totalContent := network.LibrarySize()
	cacheCapacity := network.GetCacheCapacity()
	numCacheServer := network.GetNumberOfCacheServers()
	numChromosome := 2
	// generation := 2
	mutanRate := float64(0.5)

	parent := initializePopulation(network, totalContent, cacheCapacity, numChromosome, numCacheServer)
	offspring := initializePopulation(network, totalContent, cacheCapacity, numChromosome, numCacheServer)

	fmt.Println("Test initializePopulation")
	fmt.Println(parent)
	fmt.Println(offspring)

	fmt.Println("Test computeFitness")
	parentFitness := computeFitness(network, parent)
	fmt.Println(parentFitness)
	fmt.Println(computeFitness(network, parent))

	fmt.Println("Test crossover")
	crossover(offspring, parent)
	fmt.Println("parent")
	fmt.Println(parent)
	fmt.Println("offspring")
	fmt.Println(offspring)

	fmt.Println("Test mutation")
	mutation(network, offspring, mutanRate, cacheCapacity, totalContent)
	fmt.Println("offspring")
	fmt.Println(offspring)

	fmt.Println("Test selection")
	offspringFitness := computeFitness(network, offspring)
	fmt.Println("Before")
	fmt.Println("offspring parent")
	fmt.Println(offspring)
	fmt.Println(parent)
	fmt.Println("offspring parent fitness")
	fmt.Println(offspringFitness)
	fmt.Println(parentFitness)

	selection(offspring, parent, offspringFitness, parentFitness)
	fmt.Println("After")
	fmt.Println("offspring parent")
	fmt.Println(offspring)
	fmt.Println(parent)
	fmt.Println("offspring parent fitness")
	fmt.Println(offspringFitness)
	fmt.Println(parentFitness)
}

/*
func main() {
	if parser.Options.GraphFilename == "" {
		utils.DebugPrint(fmt.Sprintln("Please specify graph filename."))
		parser.PrintDefaults()
		os.Exit(1)
	}

	network, _ := graphLoader.LoadGraph(parser.Options.GraphFilename)
	// warmupRequestCount := parser.Options.WarmupRequestCount
	// evaluationRequestCount := parser.Options.EvaluationRequestCount

	totalContent := network.LibrarySize()
	cacheCapacity := network.GetCacheCapacity()
	numCacheServer := network.GetNumberOfCacheServers()
	numChromosome := 10
	generation := 1000000
	mutanRate := float64(0.005)

	parent := initializePopulation(totalContent, cacheCapacity, numChromosome, numCacheServer)
	offspring := initializePopulation(totalContent, cacheCapacity, numChromosome, numCacheServer)

	parentFitness := computeFitness(network, parent)
	//offspringFitness := computeFitness(network, offspring)
	var offspringFitness []float64

	var temp float64

	start := time.Now()
	for i := 0; i < generation; i++ {
		crossover(offspring, parent)
		offspringFitness = computeFitness(network, offspring)
		selection(offspring, parent, offspringFitness, parentFitness)
		mutation(offspring, mutanRate, cacheCapacity, totalContent)
		offspringFitness = computeFitness(network, offspring)
		selection(offspring, parent, offspringFitness, parentFitness)

		if i == 0 || temp != offspringFitness[0] {
			fmt.Printf("Generation %d: %f\n", i, offspringFitness[0])
		}

		temp = offspringFitness[0]
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("%s ", "Elapsed time GA(s): ")
	fmt.Println(elapsed)
	// gaTest(network)

}
*/
