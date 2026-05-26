package engine

type bucket struct {
	counts map[string]int64
}

func newBucket() *bucket {
	return &bucket{
		counts: make(map[string]int64),
	}
}

func (b *bucket) add(query string) {
	b.counts[query]++
}
