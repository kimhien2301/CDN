package loader

import (
	"cache/eviction/iris"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
)

func LookupPlugin(filename string) (string, bool) {
	candidates := make([]string, 0)
	execpath, err := os.Executable()
	execdir := filepath.Dir(execpath)
	if err == nil {
		candidates = append(candidates, filepath.Join(execdir, filename))
		candidates = append(candidates, filepath.Join(execdir, "plugin", filename))
	}
	candidates = append(candidates, filepath.Join(".", filename))
	candidates = append(candidates, filepath.Join(".", "plugin", filename))
	candidates = append(candidates, filepath.Join(os.Getenv("GOPATH"), "bin", "plugin", filename))

	for _, candidate := range candidates {
		_, err = os.Stat(candidate)
		if err == nil {
			return candidate, true
		}
	}

	return "", false
}

func LoadPlugin(capacity int, spectrumRatio float64) (iris.Accessor, bool) {
	filename, exist := LookupPlugin("graph.so")
	if !exist {
		fmt.Println("Failed to find plugin: graph.so")
		return nil, false
	}

	p, err := plugin.Open(filename)
	if err != nil {
		fmt.Printf("Failed to load plugin: %s\n", filename)
		return nil, false
	}

	newfuncsym, err := p.Lookup("NewIrisCache")
	if err != nil {
		fmt.Println("Failed to lookup symbol: NewIrisCache()")
		return nil, false
	}

	newfunc, ok := newfuncsym.(func(int, float64) iris.Accessor)
	if !ok {
		fmt.Println("Failed for type assertion: iris.Accessor")
		return nil, false
	}

	accessor := newfunc(capacity, spectrumRatio)

	return accessor, true
}
