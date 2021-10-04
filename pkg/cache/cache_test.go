package cache_test

import (
	"testing"

	"scroll/pkg/cache"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {

	assert := assert.New(t)

	lru := cache.NewLru(4)

	lru.Put(1, 2)

	v := lru.Get(1)

	assert.Equal(2, v)
}

func TestFull(t *testing.T) {

	assert := assert.New(t)

	maxSize := 4

	lru := cache.NewLru(maxSize)

	for i := 0; i < maxSize; i++ {
		lru.Put(i, 10*i)
	}

	for i := 0; i < maxSize; i++ {
		v := lru.Get(i)

		assert.Equal(10*i, v)
	}

	stats := lru.GetStats()
	assert.Equal(maxSize, stats.Hits)
	assert.Equal(0, stats.Misses)
	assert.Equal(0, stats.Evictions)
}

func TestEvictOne(t *testing.T) {

	assert := assert.New(t)

	maxSize := 4

	lru := cache.NewLru(maxSize)

	for i := 0; i < maxSize+1; i++ {
		lru.Put(i, 10*i)
	}

	assert.Equal(-1, lru.Get(0))
	assert.Equal(10*4, lru.Get(4))

	stats := lru.GetStats()
	assert.Equal(1, stats.Hits)
	assert.Equal(1, stats.Misses)
	assert.Equal(1, stats.Evictions)
}

func TestWikipediaValues(t *testing.T) {

	assert := assert.New(t)

	// https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU)
	accessSequence := []int{'A', 'B', 'C', 'D', 'E', 'D', 'F'}
	cachedKeys := []int{'E', 'F', 'C', 'D'}
	evictedKeys := []int{'A', 'B'}

	lru := cache.NewLru(4)

	for _, k := range accessSequence {
		if lru.Get(k) == -1 {
			lru.Put(k, k*10)
		}
	}

	stats := lru.GetStats()
	assert.Equal(1, stats.Hits)
	assert.Equal(6, stats.Misses)
	assert.Equal(2, stats.Evictions)

	for _, k := range cachedKeys {
		v := lru.Get(k)
		assert.NotEqual(-1, v, "for key %c", k)
		assert.Equal(10*k, v, "for key %c", k)
	}

	for _, k := range evictedKeys {
		assert.Equal(-1, lru.Get(k), "key %c found", k)
	}
}

func BenchmarkNoEvict(b *testing.B) {

	lruSize := 16 * 1024
	lru := cache.NewLru(lruSize)

	for i := 0; i < lruSize; i++ {
		lru.Put(i, i)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < lruSize; i++ {
			lru.Put(i, i)
		}
	}
	b.StopTimer()

	stats := lru.GetStats()
	assert.Equal(b, 0, stats.Evictions)
}

func BenchmarkEvict(b *testing.B) {

	lruSize := 16 * 1024
	lru := cache.NewLru(lruSize)

	for i := 0; i < lruSize; i++ {
		lru.Put(i, i)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < lruSize; i++ {
			lru.Put(i+lruSize, i+lruSize)
		}
	}
	b.StopTimer()

	stats := lru.GetStats()
	assert.Equal(b, lruSize, stats.Evictions)
}
