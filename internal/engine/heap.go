package engine

import "container/heap"

type queryScore struct {
	query string
	count int64
}

type minHeap []queryScore

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].count < h[j].count }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)        { *h = append(*h, x.(queryScore)) }
func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func topN(counts map[string]int64, n int) []queryScore {
	h := &minHeap{}
	heap.Init(h)

	for query, count := range counts {
		if h.Len() < n {
			heap.Push(h, queryScore{query, count})
		} else if count > (*h)[0].count {
			(*h)[0] = queryScore{query, count}
			heap.Fix(h, 0)
		}
	}

	result := make([]queryScore, h.Len())
	for i := h.Len() - 1; i >= 0; i-- {
		result[i] = heap.Pop(h).(queryScore)
	}
	return result
}
