package dags

import (
	"shanhu.io/std/graph"
)

// FromStdGraph converts a std/graph.Graph into a dags.Graph. Every node in g
// becomes a key in the result, and each edge From->To adds To to the outgoing
// list of From.
func FromStdGraph(g *graph.Graph) *Graph {
	nodes := make(map[string][]string)
	for _, n := range g.Nodes {
		nodes[n.Name] = nil // touch every node first
	}
	for _, e := range g.Edges {
		nodes[e.From] = append(nodes[e.From], e.To)
	}
	return &Graph{Nodes: nodes}
}
