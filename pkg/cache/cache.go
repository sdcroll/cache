package cache

import (
	"container/list"
	"log"
)

type Stats struct {
	Hits      int
	Misses    int
	Evictions int
}

type Entry struct {
	key   int
	value int
}

type Lru struct {
	lruMap  map[int]*list.Element // map key -> element on the lruList
	lruList *list.List            // list of Entry ptrs, most recently accessed at front of list
	maxSize int
	stats   Stats
}

func NewLru(maxSize int) *Lru {

	if maxSize <= 0 {
		log.Fatalf("maxSize must be positive")
	}

	return &Lru{
		lruMap:  make(map[int]*list.Element),
		lruList: list.New(),
		maxSize: maxSize,
	}
}

func (l *Lru) Put(key, value int) {

	element, ok := l.lruMap[key]

	// if key exists, move element to front of list (now most recently accessed)

	if ok {
		l.lruList.MoveToFront(element)
		return
	}

	// key not found, add new entry

	l.lruMap[key] = l.lruList.PushFront(&Entry{key, value})

	// return if space available

	if len(l.lruMap) <= l.maxSize {
		return
	}

	// exceeded max size, evict least recently used (found at back of list)

	ev := l.lruList.Remove(l.lruList.Back())
	delete(l.lruMap, ev.(*Entry).key)

	l.stats.Evictions += 1
}

func (l *Lru) Get(key int) int {

	element, ok := l.lruMap[key]

	// return failure if not found in cache

	if !ok {
		l.stats.Misses += 1
		return -1
	}

	// move element to front of list (now most recently accessed)

	l.lruList.MoveToFront(element)

	l.stats.Hits += 1
	return element.Value.(*Entry).value
}

func (l *Lru) GetStats() Stats {
	return l.stats
}
