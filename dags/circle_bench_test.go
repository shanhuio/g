package dags

import (
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkFindCircle(b *testing.B) {
	for b.Loop() {
		benchmarkFindCircle()
	}
}

func benchmarkFindCircle() []*MapNode {
	ret := make(map[string][]string)

	for i := range 1000 {
		var edge []string
		for j := range 1000 {
			if i == j {
				continue
			}
			q := rand.Int31n(500)
			if q < 1 {
				var flag bool
				for _, ele := range ret[strconv.Itoa(j)] {
					if ele == strconv.Itoa(i) {
						flag = true
					}
				}
				if !flag {
					edge = append(edge, strconv.Itoa(j))
				}
			}
		}
		ret[strconv.Itoa(i)] = edge
	}

	g := &Graph{Nodes: ret}
	nodes, _ := initMap(g)

	return minCircle(nodes.Nodes)
}
