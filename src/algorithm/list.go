package algorithm

type Entry struct {
	Key   interface{}
	Value float64
}

type List []Entry

func (list List) Len() int {
	return len(list)
}

func (list List) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list List) Less(i, j int) bool {
	return (list[i].Value < list[j].Value)
}
