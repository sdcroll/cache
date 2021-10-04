package cache_test

import (
	"testing"

	"scroll/pkg/cache"
)

func TestBasic(t *testing.T) {

	lru := cache.NewLru(4)

	lru.Put(1, 2)

	v := lru.Get(1)

	if v != 2 {
		t.Errorf("expected 2, got %v", v)
	}
}

func TestFull(t *testing.T) {

	maxSize := 4

	lru := cache.NewLru(maxSize)

	for i := 0; i < maxSize; i++ {
		lru.Put(i, i)
	}

	for i := 0; i < maxSize; i++ {
		v := lru.Get(i)

		if v != i {
			t.Errorf("expected %v, got %v", i, v)
		}
	}
}

func TestEvictOne(t *testing.T) {

	maxSize := 4

	lru := cache.NewLru(maxSize)

	for i := 0; i < maxSize+1; i++ {
		lru.Put(i, i)
	}

	v := lru.Get(0)

	if v != -1 {
		t.Errorf("expected %v, got %v", -1, v)
	}

	v = lru.Get(4)

	if v != 4 {
		t.Errorf("expected %v, got %v", 4, v)
	}
}
