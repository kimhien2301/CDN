package userdist

import "fmt"
import "io/ioutil"
import "strings"
import "strconv"
import "sort"
import "math/rand"
import "algorithm"

type Distribution struct {
	rawPopularities algorithm.List
	pdfMemo         []float64
	cdfMemo         []float64
}

func New(filename string) *Distribution {
	dist := new(Distribution)
	dist.loadFile(filename)
	return dist
}

func (dist *Distribution) Len() int {
	return len(dist.pdfMemo)
}

func (dist *Distribution) validatePopularity() {
	dist.pdfMemo = make([]float64, 0)
	dist.cdfMemo = make([]float64, 0)

	sum := 0.0
	for _, entry := range dist.rawPopularities {
		sum += entry.Value
	}

	dist.pdfMemo = append(dist.pdfMemo, dist.rawPopularities[0].Value/sum)
	dist.cdfMemo = append(dist.cdfMemo, dist.rawPopularities[0].Value/sum)
	for index := 1; index < len(dist.rawPopularities); index++ {
		dist.pdfMemo = append(dist.pdfMemo, dist.rawPopularities[index].Value/sum)
		dist.cdfMemo = append(dist.cdfMemo, dist.cdfMemo[index-1]+dist.rawPopularities[index].Value/sum)
	}
}

func (dist *Distribution) loadFile(filename string) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("error in reading file: %s", err))
	}
	lines := strings.Split(string(contents), "\n")
	for index := range lines {
		line := lines[index]
		popularity, err := strconv.ParseFloat(line, 64)
		if err != nil {
			continue
		}
		dist.rawPopularities = append(dist.rawPopularities, algorithm.Entry{line, -popularity})
	}
	sort.Sort(dist.rawPopularities)
	dist.validatePopularity()
}

func (dist *Distribution) PDF(rank int) float64 {
	return dist.pdfMemo[rank-1]
}

func (dist *Distribution) CDF(rank int) float64 {
	return dist.cdfMemo[rank-1]
}

func (dist *Distribution) Intn() int {
	var rank int
	mark := rand.Float64()
	for rank = 1; rank <= dist.Len(); rank++ {
		if dist.cdfMemo[rank-1] > mark {
			break
		}
	}
	return rank
}
