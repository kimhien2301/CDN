package algorithm

func Combination(src []interface{}, pick int) [][]interface{} {
	combinations := make([][]interface{}, 0)
	if pick <= 0 {
		return combinations
	} else if pick == 1 {
		for _, entry := range src {
			combination := make([]interface{}, 0)
			combination = append(combination, entry)
			combinations = append(combinations, combination)
		}
		return combinations
	}

	for index := range src {
		if index > len(src)-pick {
			break
		}
		for _, combination := range Combination(src[index+1:len(src)], pick-1) {
			concat := make([]interface{}, 0)
			concat = append(concat, src[index])
			for _, entry := range combination {
				concat = append(concat, entry)
			}
			combinations = append(combinations, concat)
		}
	}

	return combinations
}
