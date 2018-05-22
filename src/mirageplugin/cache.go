package main

import (
	"cache/eviction/iris"
	"cache/eviction/lfu"
	"cache/eviction/modifiedlru"
	"fmt"
)

type CacheStorage struct {
	// spectrumCache    *modifiedlru.CacheStorage
	spectrumCache *lfu.CacheStorage // add
	mlruCache     *modifiedlru.CacheStorage
	// arcCache         *arc.CacheStorage // add
	serverSpectrum   uint64
	contentSpectrums []uint64
	spectrumCapacity int
	mlruCapacity     int
	// arcCapacity      int // add
	hitCount     int
	missCount    int
	spectrumTags map[int][]uint64
}

func NewIrisCache(capacity int, spectrumRatio float64) iris.Accessor {
	irisCache := new(CacheStorage)
	irisCache.hitCount = 0
	irisCache.missCount = 0
	irisCache.spectrumCapacity = int(float64(capacity) * spectrumRatio)
	irisCache.mlruCapacity = capacity - irisCache.spectrumCapacity
	// irisCache.arcCapacity = capacity - irisCache.spectrumCapacity // add
	irisCache.spectrumCache = lfu.New(irisCache.spectrumCapacity)
	irisCache.mlruCache = modifiedlru.New(irisCache.mlruCapacity, 5) // jump = 5
	// irisCache.arcCache = arc.New(irisCache.arcCapacity) // add
	// fmt.Printf("Spectrum Capacity: %d\nModified LRU Capacity: %d\n", irisCache.spectrumCapacity, irisCache.mlruCapacity)
	return irisCache
}

func (cache *CacheStorage) CacheList() []interface{} {
	cacheList := make([]interface{}, 0)
	return cacheList
}

func (cache *CacheStorage) Inspect() {
}

func (cache *CacheStorage) FillUp() {
}

func (cache *CacheStorage) SetContentSpectrums(contentSpectrums []uint64) {
	cache.contentSpectrums = contentSpectrums
}

func (cache *CacheStorage) SetServerSpectrum(spectrum uint64) {
	cache.serverSpectrum = spectrum
}

func (cache *CacheStorage) Len() int {
	return cache.spectrumCache.Len() + cache.mlruCache.Len()
	// return cache.spectrumCache.Len() + cache.arcCache.Len() // add
}

func (cache *CacheStorage) Exist(key interface{}) bool {
	return cache.spectrumCache.Exist(key) || cache.mlruCache.Exist(key)
	// return cache.spectrumCache.Exist(key) || cache.arcCache.Exist(key) // add
}

////////////////////////////////////////////////////////////////

func matchWithServerColor(contentTag []uint64, serverSpectrum uint64) bool {
	if contentTag[serverSpectrum] == 1 {
		return true
	}
	return false
}

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	// fmt.Println(key)
	// fmt.Println(cache.spectrumTags[key.(int)])
	var hasTag = true
	if len(cache.spectrumTags[key.(int)]) == 0 {
		hasTag = false
	}
	if hasTag && matchWithServerColor(cache.spectrumTags[key.(int)], cache.serverSpectrum) {
		// fmt.Println("Spectrum Caches 1")
		if cache.spectrumCache.Len() >= cache.spectrumCache.Capacity() {
			spectrumCacheList := cache.spectrumCache.CacheList()
			// fmt.Println("spectrumCacheList")
			// fmt.Println(spectrumCacheList)
			for _, spectrumCacheListKey := range spectrumCacheList {
				if !matchWithServerColor(cache.spectrumTags[spectrumCacheListKey.(int)], cache.serverSpectrum) {
					// fmt.Println("EvictWithKey")
					cache.spectrumCache.EvictWithKey(spectrumCacheListKey)
					break
				}
			}
		}
		cache.spectrumCache.Insert(key, value)
		// fmt.Println(:Hybrid Caches 2")
		return value
	} else {
		// fmt.Println("lru Caches 1")
		cache.mlruCache.Insert(key, value)
		// cache.arcCache.Insert(key, value) // add
		return value
	}
	return nil
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	data := cache.spectrumCache.Fetch(key)
	if data != nil {
		cache.hitCount++
		fmt.Println("HIT")
		return data
	}

	data = cache.mlruCache.Fetch(key)
	// data = cache.arcCache.Fetch(key)
	if data != nil {
		cache.hitCount++
		fmt.Println("HIT")
		return data
	}

	cache.missCount++
	fmt.Println("MISS")
	return nil
}

/////////////////////////////////////////////////////////////////

func (cache *CacheStorage) SpectrumCapacity() int {
	return cache.spectrumCapacity
}

func (cache *CacheStorage) Capacity() int {
	return cache.spectrumCapacity + cache.mlruCapacity
	// return cache.spectrumCapacity + cache.arcCapacity // add
}

func (cache *CacheStorage) HitCount() int {
	return cache.hitCount
}

func (cache *CacheStorage) MissCount() int {
	return cache.missCount
}

func (cache *CacheStorage) ResetCount() {
	cache.hitCount = 0
	cache.missCount = 0
}

func (cache *CacheStorage) Clear() {
	cache.ResetCount()
	cache.spectrumCache.Clear()
	cache.mlruCache.Clear()
	// cache.arcCache.Clear() // add
}
