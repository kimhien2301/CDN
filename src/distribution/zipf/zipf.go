package zipf

import "math"
import "math/rand"

type Distribution struct {
	skewness    float64
	length      int
	denominator float64
	pdfMemo     []float64
	cdfMemo     []float64
}

func New(skewness float64, length int) *Distribution {
	zipf := new(Distribution)
	zipf.skewness = skewness
	zipf.length = length
	zipf.Init()
	return zipf
}

func (zipf *Distribution) Init() {
	zipf.pdfMemo = make([]float64, 0)
	zipf.cdfMemo = make([]float64, 0)

	zipf.denominator = 0.0
	for index := 1; index <= zipf.length; index++ {
		zipf.denominator += 1.0 / math.Pow(float64(index), zipf.skewness)
	}

	zipf.pdfMemo = append(zipf.pdfMemo, (1.0/math.Pow(1.0, zipf.skewness))/zipf.denominator)
	zipf.cdfMemo = append(zipf.pdfMemo, (1.0/math.Pow(1.0, zipf.skewness))/zipf.denominator)
	for rank := 2; rank <= zipf.length; rank++ {
		zipf.pdfMemo = append(zipf.pdfMemo, (1.0/math.Pow(float64(rank), zipf.skewness))/zipf.denominator)
		zipf.cdfMemo = append(zipf.cdfMemo, zipf.cdfMemo[rank-2]+(1.0/math.Pow(float64(rank), zipf.skewness))/zipf.denominator)
	}
}

func (zipf *Distribution) PDF(rank int) float64 {
	return zipf.pdfMemo[rank-1]
}

func (zipf *Distribution) CDF(rank int) float64 {
	return zipf.cdfMemo[rank-1]
}

func (zipf *Distribution) Intn() int {
	if zipf.denominator == 0.0 {
		zipf.Init()
	}
	var rank int
	mark := rand.Float64()
	for rank = 1; rank <= zipf.length; rank++ {
		if zipf.cdfMemo[rank-1] > mark {
			break
		}
	}
	return rank
}
