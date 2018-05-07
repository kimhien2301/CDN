package gamma

import "math"
import "math/rand"

type Distribution struct {
	K, Theta float64
	Length   int
	sum      float64
	pdfMemo  []float64
	cdfMemo  []float64
}

func New(k, theta float64, length int) *Distribution {
	gamma := new(Distribution)
	gamma.K = k
	gamma.Theta = theta
	gamma.Length = length
	return gamma
}

func (dist *Distribution) PDF(rank int) float64 {
	if dist.sum == 0.0 {
		dist.Init()
	}
	return dist.pdfMemo[rank-1]
}

func (dist *Distribution) CDF(rank int) float64 {
	if dist.sum == 0.0 {
		dist.Init()
	}
	return dist.cdfMemo[rank-1]
}

func (dist *Distribution) Intn() int {
	if dist.sum == 0.0 {
		dist.Init()
	}
	var rank int
	mark := rand.Float64()
	for rank = 1; rank <= dist.Length; rank++ {
		if dist.cdfMemo[rank-1] > mark {
			break
		}
	}
	return rank
}

func (dist *Distribution) pdf(x float64) float64 {
	return math.Pow(x, dist.K-1) * math.Exp(-x/dist.Theta) / (math.Gamma(dist.K) * math.Pow(dist.Theta, dist.K))
}

func (dist *Distribution) Init() {
	dist.pdfMemo = make([]float64, 0)
	dist.cdfMemo = make([]float64, 0)

	dist.sum = 0.0
	for index := 1; index <= dist.Length; index++ {
		dist.sum += dist.pdf(float64(index))
	}

	dist.pdfMemo = append(dist.pdfMemo, dist.pdf(1.0)/dist.sum)
	dist.cdfMemo = append(dist.cdfMemo, dist.pdf(1.0)/dist.sum)
	for index := 2; index <= dist.Length; index++ {
		dist.pdfMemo = append(dist.pdfMemo, dist.pdf(float64(index))/dist.sum)
		dist.cdfMemo = append(dist.cdfMemo, dist.cdfMemo[index-2]+dist.pdf(float64(index))/dist.sum)
	}
}
