package dags

import (
	"testing"

	"shanhu.io/std/graph"
)

func TestFromStdGraph(t *testing.T) {
	g := &graph.Graph{
		Nodes: []*graph.Node{
			{Name: "a", Comment: "node a"},
			{Name: "b"},
			{Name: "c"},
		},
		Edges: []*graph.Edge{
			{From: "a", To: "b"},
			{From: "a", To: "c"},
			{From: "b", To: "c"},
		},
	}

	got := FromStdGraph(g)
	checkGraphEqual(t, got, map[string][]string{
		"a": {"b", "c"},
		"b": {"c"},
		"c": nil,
	})
}

func TestFromStdGraphEmpty(t *testing.T) {
	got := FromStdGraph(&graph.Graph{})
	checkGraphEqual(t, got, map[string][]string{})
}

func TestFromStdGraphIsolatedNodes(t *testing.T) {
	g := &graph.Graph{
		Nodes: []*graph.Node{{Name: "a"}, {Name: "b"}},
	}

	got := FromStdGraph(g)
	checkGraphEqual(t, got, map[string][]string{
		"a": nil,
		"b": nil,
	})
}

func TestFromStdGraphEdgeOrder(t *testing.T) {
	// Outgoing targets preserve edge order, not sorted order.
	g := &graph.Graph{
		Nodes: []*graph.Node{{Name: "a"}, {Name: "b"}, {Name: "c"}},
		Edges: []*graph.Edge{
			{From: "a", To: "c"},
			{From: "a", To: "b"},
		},
	}

	got := FromStdGraph(g)
	checkGraphEqual(t, got, map[string][]string{
		"a": {"c", "b"},
		"b": nil,
		"c": nil,
	})
}
