package utils

import "fmt"
import "io/ioutil"
import "encoding/json"

type SeparatorRanks []int

type BestSeparatorRanks struct {
	NetworkID          string         `json:"NetworkID"`
	RequestModel       string         `json:"RequestModel"`
	RequestModelParams []float64      `json:"RequestModelParams"`
	Ranks              SeparatorRanks `json:"SeparatorRanks"`
}

type MirageStore struct {
	ReferenceSeparatorRanks []SeparatorRanks     `json:"ReferenceSeparatorRanks"`
	BestSeparatorRanks      []BestSeparatorRanks `json:"BestSeparatorRanks"`
}

func NewMirageStore() MirageStore {
	var mirageStore MirageStore
	mirageStore.ReferenceSeparatorRanks = append(mirageStore.ReferenceSeparatorRanks, []int{36, 42, 54, 268}) // 0.1
	mirageStore.ReferenceSeparatorRanks = append(mirageStore.ReferenceSeparatorRanks, []int{32, 36, 54, 278}) // 0.2
	mirageStore.ReferenceSeparatorRanks = append(mirageStore.ReferenceSeparatorRanks, []int{28, 32, 52, 288}) // 0.3
	mirageStore.ReferenceSeparatorRanks = append(mirageStore.ReferenceSeparatorRanks, []int{24, 27, 48, 301}) // 0.4
	mirageStore.ReferenceSeparatorRanks = append(mirageStore.ReferenceSeparatorRanks, []int{20, 28, 40, 312}) // 0.5
	mirageStore.ReferenceSeparatorRanks = append(mirageStore.ReferenceSeparatorRanks, []int{16, 20, 32, 332}) // 0.6
	mirageStore.ReferenceSeparatorRanks = append(mirageStore.ReferenceSeparatorRanks, []int{8, 16, 28, 348})  // 0.7
	mirageStore.ReferenceSeparatorRanks = append(mirageStore.ReferenceSeparatorRanks, []int{4, 4, 20, 372})   // 0.8
	mirageStore.ReferenceSeparatorRanks = append(mirageStore.ReferenceSeparatorRanks, []int{0, 0, 12, 388})   // 0.9
	return mirageStore
}

func StoreData(filename string, mirageStore MirageStore) {
	bytes, _ := json.Marshal(mirageStore)
	ioutil.WriteFile(filename, bytes, 0644)
}

func LoadData(filename string) MirageStore {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return NewMirageStore()
	}

	var mirageStore MirageStore
	if err = json.Unmarshal(bytes, &mirageStore); err != nil {
		panic(fmt.Errorf("Fatal error in Unmarshal(): %s", err))
	}

	return mirageStore
}
