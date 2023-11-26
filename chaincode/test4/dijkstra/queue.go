package dijkstra

import "sort"

// Queue is a basic priority queue implementation, where the Node with the
// lowest priority is kept as first element in the queue
type Queue struct {
	Keys  []string
	Nodes map[string]int
}

// Len is part of sort.Interface
func (q *Queue) Len() int {
	return len(q.Keys)
}

// Swap is part of sort.Interface
func (q *Queue) Swap(i, j int) {
	q.Keys[i], q.Keys[j] = q.Keys[j], q.Keys[i]
}

// Less is part of sort.Interface
func (q *Queue) Less(i, j int) bool {
	a := q.Keys[i]
	b := q.Keys[j]

	return q.Nodes[a] < q.Nodes[b]
}

// Set updates or inserts a new Key in the priority queue
func (q *Queue) Set(Key string, priority int) {
	// inserts a new Key if we don't have it already
	if _, ok := q.Nodes[Key]; !ok {
		q.Keys = append(q.Keys, Key)
	}

	// set the priority for the Key
	q.Nodes[Key] = priority

	// sort the Keys array
	sort.Sort(q)
}

// Next removes the first element from the queue and retuns it's Key and priority
func (q *Queue) Next() (Key string, priority int) {
	// shift the Key form the queue
	Key, Keys := q.Keys[0], q.Keys[1:]
	q.Keys = Keys

	priority = q.Nodes[Key]

	delete(q.Nodes, Key)

	return Key, priority
}

// IsEmpty returns true when the queue is empty
func (q *Queue) IsEmpty() bool {
	return len(q.Keys) == 0
}

// Get returns the priority of a passed Key
func (q *Queue) Get(Key string) (priority int, ok bool) {
	priority, ok = q.Nodes[Key]
	return
}

// NewQueue creates a new empty priority queue
func NewQueue() *Queue {
	var q Queue
	q.Nodes = make(map[string]int)
	return &q
}
