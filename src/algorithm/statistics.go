package algorithm

import "math"

func Sum(slice []float64) float64 {
	var sum float64
	for _, value := range slice {
		sum += value
	}
	return sum
}

func Average(slice []float64) float64 {
	return Sum(slice) / float64(len(slice))
}

func Variance(slice []float64) float64 {
	variance := 0.0
	average := Average(slice)
	for _, value := range slice {
		diff := value - average
		variance += diff * diff
	}
	return variance
}

func StandardDeviation(slice []float64) float64 {
	return math.Sqrt(Variance(slice))
}

func RSS(measured, ideal []float64) float64 {
	rss := 0.0
	for index := range measured {
		diff := measured[index] - ideal[index]
		rss += diff * diff
	}
	return rss
}
