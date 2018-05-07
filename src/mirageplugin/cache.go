package main

import (
	"cache/eviction/iris"
	"cache/eviction/lfu"
	"cache/eviction/modifiedlru"
)

type CacheStorage struct {
	//spectrumCache    *modifiedlru.CacheStorage
	spectrumCache    *lfu.CacheStorage //
	mlruCache        *modifiedlru.CacheStorage
	serverSpectrum   uint64
	contentSpectrums []uint64
	spectrumCapacity int
	mlruCapacity     int
	hitCount         int
	missCount        int
	spectrumTags     map[int][]uint64
}

func NewIrisCache(capacity int, spectrumRatio float64) iris.Accessor {
	irisCache := new(CacheStorage)
	irisCache.hitCount = 0
	irisCache.missCount = 0
	irisCache.spectrumCapacity = int(float64(capacity) * spectrumRatio)
	irisCache.mlruCapacity = capacity - irisCache.spectrumCapacity
	irisCache.spectrumCache = lfu.New(irisCache.spectrumCapacity)
	irisCache.mlruCache = modifiedlru.New(irisCache.mlruCapacity, 5) // jump = 5
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
}

func (cache *CacheStorage) Exist(key interface{}) bool {
	return cache.spectrumCache.Exist(key) || cache.mlruCache.Exist(key)
}

////////////////////////////////////////////////////////////////

func matchWithServerColor(contentTag []uint64, serverSpectrum uint64) bool {
	if contentTag[serverSpectrum] == 1 {
		return true
	}
	return false
}

func (cache *CacheStorage) Insert(key, value interface{}) interface{} {
	var hasTag = true
	if len(cache.spectrumTags[key.(int)]) == 0 {
		hasTag = false
	}
	if hasTag && matchWithServerColor(cache.spectrumTags[key.(int)], cache.serverSpectrum) {
		if cache.spectrumCache.Len() >= cache.spectrumCache.Capacity() {
			spectrumCacheList := cache.spectrumCache.CacheList()
			for _, spectrumCacheListKey := range spectrumCacheList {
				if !matchWithServerColor(cache.spectrumTags[spectrumCacheListKey.(int)], cache.serverSpectrum) {
					cache.spectrumCache.EvictWithKey(spectrumCacheListKey)
					break
				}
			}
		}
		cache.spectrumCache.Insert(key, value)
		return value
	} else {
		cache.mlruCache.Insert(key, value)
		return value
	}
	return nil
}

func (cache *CacheStorage) Fetch(key interface{}) interface{} {
	data := cache.spectrumCache.Fetch(key)
	if data != nil {
		cache.hitCount++
		return data
	}

	data = cache.mlruCache.Fetch(key)
	if data != nil {
		cache.hitCount++
		return data
	}

	cache.missCount++
	return nil
}

/////////////////////////////////////////////////////////////////

func (cache *CacheStorage) SpectrumCapacity() int {
	return cache.spectrumCapacity
}

func (cache *CacheStorage) Capacity() int {
	return cache.spectrumCapacity + cache.mlruCapacity
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
}
