package engine

import (
	"context"
	"sync"
	"time"
)

const bucketCount = 300

type Entry struct {
	Query string `json:"query"`
	Count int64  `json:"count"`
}

type Engine struct {
	mu      sync.RWMutex
	buckets [bucketCount]*bucket
	current int // newest bucket index

	stoplist map[string]struct{}
	slMu     sync.RWMutex

	cachedTop []Entry
	cachedMu  sync.RWMutex
}

func NewEngine() *Engine {
	e := &Engine{
		stoplist: make(map[string]struct{}),
	}

	for i := range e.buckets {
		e.buckets[i] = newBucket()
	}

	return e
}

func (e *Engine) Add(query string) {
	e.mu.Lock()
	e.buckets[e.current].add(query)
	e.mu.Unlock()
}

func (e *Engine) Top(n int) []Entry {
	e.cachedMu.RLock()
	defer e.cachedMu.RUnlock()

	if n > len(e.cachedTop) {
		n = len(e.cachedTop)
	}

	result := make([]Entry, n)
	copy(result, e.cachedTop[:n])
	return result
}

func (e *Engine) AddToStoplist(word string) {
	e.slMu.Lock()
	e.stoplist[word] = struct{}{}
	e.slMu.Unlock()
}

func (e *Engine) RemoveFromStoplist(word string) {
	e.slMu.Lock()
	delete(e.stoplist, word)
	e.slMu.Unlock()
}

func (e *Engine) IsBlocked(query string) bool {
	e.slMu.RLock()
	_, blocked := e.stoplist[query]
	e.slMu.RUnlock()
	return blocked
}

func (e *Engine) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.rotate()
			e.recalculate()
		}
	}
}

func (e *Engine) rotate() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.current = (e.current + 1) % bucketCount
	e.buckets[e.current] = newBucket()
}

func (e *Engine) recalculate() {
	e.mu.RLock()

	merged := make(map[string]int64)
	for _, b := range e.buckets {
		for q, c := range b.counts {
			merged[q] += c
		}
	}

	e.mu.RUnlock()

	e.slMu.RLock()
	for word := range e.stoplist {
		delete(merged, word)
	}
	e.slMu.RUnlock()

	scores := topN(merged, 100)
	top := make([]Entry, len(scores))
	for i, s := range scores {
		top[i] = Entry{Query: s.query, Count: s.count}
	}

	e.cachedMu.Lock()
	e.cachedTop = top
	e.cachedMu.Unlock()
}
