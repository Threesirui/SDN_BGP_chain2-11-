/*
Package dijkstra is an highly optimised implementation of the Dijkstra
algorithm, used for find the shortest path between points of a graph.

A graph is a map of points and map to the neighbouring points in the graph and
the Cost to reach them.
A trivial example of a graph definition is:

	Graph{
		"a": {"b": 10, "c": 20},
		"b": {"a": 50},
		"c": {"b": 10, "a": 25},
	}

*/
package dijkstra

import "fmt"

type Node struct {
	Key  string
	Cost int
}

// Graph is a rappresentation of how the points in our graph are connected
// between each other
type Graph map[string]map[string]int

// Path finds the shortest path between start and target, also returning the
// total Cost of the found path.
func (g Graph) Path(start, target string) (path []string, Cost int, err error) {
	if len(g) == 0 {
		err = fmt.Errorf("cannot find path in empty map")
		return
	}

	// ensure start and target are part of the graph
	if _, ok := g[start]; !ok {
		err = fmt.Errorf("cannot find start %v in graph", start)
		return
	}
	if _, ok := g[target]; !ok {
		err = fmt.Errorf("cannot find target %v in graph", target)
		return
	}

	explored := make(map[string]bool)   // set of Nodes we already explored
	frontier := NewQueue()              // queue of the Nodes to explore
	previous := make(map[string]string) // previously visited Node

	// add starting point to the frontier as it'll be the first Node visited
	frontier.Set(start, 0)

	// run until we visited every Node in the frontier
	for !frontier.IsEmpty() {
		// get the Node in the frontier with the lowest Cost (or priority)
		aKey, aPriority := frontier.Next()
		n := Node{aKey, aPriority}

		// when the Node with the lowest Cost in the frontier is target, we can
		// compute the Cost and path and exit the loop
		if n.Key == target {
			Cost = n.Cost

			nKey := n.Key
			for nKey != start {
				path = append(path, nKey)
				nKey = previous[nKey]
			}

			break
		}

		// add the current Node to the explored set
		explored[n.Key] = true

		// loop all the neighboring Nodes
		for nKey, nCost := range g[n.Key] {
			// skip alreadt-explored Nodes
			if explored[nKey] {
				continue
			}

			// if the Node is not yet in the frontier add it with the Cost
			if _, ok := frontier.Get(nKey); !ok {
				previous[nKey] = n.Key
				frontier.Set(nKey, n.Cost+nCost)
				continue
			}

			frontierCost, _ := frontier.Get(nKey)
			NodeCost := n.Cost + nCost

			// only update the Cost of this Node in the frontier when
			// it's below what's currently set
			if NodeCost < frontierCost {
				previous[nKey] = n.Key
				frontier.Set(nKey, NodeCost)
			}
		}
	}

	// add the origin at the end of the path
	path = append(path, start)

	// reverse the path because it was popilated
	// in reverse, form target to start
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return
}
