package distribution

type Distribution interface {
	PDF(rank int) float64
	CDF(rank int) float64
	Intn() int
}
