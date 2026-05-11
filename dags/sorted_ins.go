package dags

import (
	"sort"
)

// AllInsSorted returns all the input nodes of a node in layering order.
func AllInsSorted(node *MapNode) []*MapNode {
	var nodes []*MapNode
	for _, in := range node.AllIns {
		nodes = append(nodes, in)
	}
	sort.Sort(byLayer{nodes})
	return nodes
}
