package cache

import (
	"log"
	"math"
)

type Stats struct {
	Hits      int
	Misses    int
	Evictions int
}

type Entry struct {
	value int
	seq   int
}

type Lru struct {
	m       map[int]*Entry
	maxSize int
	nextSeq int
	stats   Stats
}

func NewLru(maxSize int) *Lru {

	if maxSize <= 0 {
		log.Fatalf("maxSize must be positive")
	}

	return &Lru{
		m:       make(map[int]*Entry),
		maxSize: maxSize,
	}
}

func (l *Lru) Put(key, value int) {

	entry, ok := l.m[key]

	// if key exists, update LRU sequence number

	if ok {
		entry.seq = l.nextSeq
		l.nextSeq += 1
		return
	}

	// key not found, add new entry

	l.m[key] = &Entry{
		value: value,
		seq:   l.nextSeq,
	}
	l.nextSeq += 1

	// return if space available

	if len(l.m) <= l.maxSize {
		return
	}

	// no space available, find and evict least recently used

	var lowKey int
	lowSeq := math.MaxInt
	for k, entry := range l.m {
		if entry.seq < lowSeq {
			lowSeq = entry.seq
			lowKey = k
		}
	}

	delete(l.m, lowKey)
	l.stats.Evictions += 1
}

func (l *Lru) Get(key int) int {

	entry, ok := l.m[key]

	// return failure if not in cache

	if !ok {
		l.stats.Misses += 1
		return -1
	}

	// return cached value

	entry.seq = l.nextSeq
	l.nextSeq += 1
	l.stats.Hits += 1

	return entry.value
}

func (l *Lru) GetStats() Stats {
	return l.stats
}
