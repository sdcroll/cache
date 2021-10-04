package cache

import (
	"log"
	"math"
)

type Entry struct {
	value int
	seq   int
}

type Lru struct {
	m       map[int]*Entry
	maxSize int
	nextSeq int
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

	if ok {
		entry.seq = l.nextSeq
		l.nextSeq += 1
		return
	}

	l.m[key] = &Entry{
		value: value,
		seq:   l.nextSeq,
	}
	l.nextSeq += 1

	if len(l.m) <= l.maxSize {
		return
	}

	var lowKey int
	lowSeq := math.MaxInt
	for k, entry := range l.m {
		if entry.seq < lowSeq {
			lowSeq = entry.seq
			lowKey = k
		}
	}

	delete(l.m, lowKey)
}

func (l *Lru) Get(key int) int {

	entry, ok := l.m[key]

	if !ok {
		return -1
	}

	entry.seq = l.nextSeq
	l.nextSeq += 1

	return entry.value
}
