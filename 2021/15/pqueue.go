package main

import (
	"container/heap"
)

/* This is a modified version of https://pkg.go.dev/container/heap#example-package-PriorityQueue
 * using generics
 * mostly for lulz, because the actual
 * implementation just suggests replacing `value string` with your real type
 * but this allowed us to test our generics
 * we have made some small modifications:
 * - Less() was replaced; it was defining a max heap where we need a min-heap
 * - PriorityQueue type went from []*Item to a struct containing items
 *   and a reverseLookup.  We needed to Peek for Dijkstra's alg if our item was still in
 *   PLUS it allows:
 * - A new updatePriority method, instead of the updateItem method that was implemented in the example
 *   this allows us to pass in just a value, rather htan keep track of *Items in our caller
 */

// An Item is something we manage in a priority queue.
type Item[T comparable] struct {
	value    T   // The value of the item; arbitrary.
	priority int // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue[T comparable] struct {
	items         []*Item[T]
	reverseLookup map[T]*Item[T]
}

func (pq PriorityQueue[T]) Len() int { return len(pq.items) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	// XXX: Changed from go example
	// they used a max-queue
	// we wanted a min-queue
	return pq.items[i].priority < pq.items[j].priority
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

func (pq *PriorityQueue[T]) Push(x interface{}) {
	// XXX we really should have a New function instead
	if pq.reverseLookup == nil {
		pq.reverseLookup = make(map[T]*Item[T])
	}

	n := len(pq.items)
	item := x.(*Item[T])
	item.index = n

	pq.items = append(pq.items, item)
	pq.reverseLookup[item.value] = item
}

func (pq *PriorityQueue[T]) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	pq.items = old[0 : n-1]

	delete(pq.reverseLookup, item.value)
	return item
}

func (pq *PriorityQueue[T]) updatePriority(value T, priority int) error {
	// XXX: because this is *Item, this can panic
	item := pq.reverseLookup[value]
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
	return nil
}

func (pq *PriorityQueue[T]) contains(value T) bool {
	_, ok := pq.reverseLookup[value]
	return ok
}
