package chain

import (
	"fmt"
	"math/rand"
)

type Node struct {
	ID    string  `json:"id"`
	Label string  `json:"label"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Size  int     `json:"size"`
}

type Edge struct {
	ID     string  `json:"id"`
	Label  string  `json:"label"`
	Size   float64 `json:"size"`
	Source string  `json:"source"`
	Target string  `json:"target"`
	Type   string  `json:"type"`
	Color  string  `json:"color"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

func NewGraph(tm [][]float64) Graph {
	nodes := make([]Node, 0, len(tm))
	for i := 0; i < len(tm); i++ {
		id := fmt.Sprintf("node-%d", i)
		nodes = append(nodes, Node{
			ID:    id,
			Label: id,
			X:     rand.Float64(),
			Y:     rand.Float64(),
			Size:  5,
		})
	}

	edges := make([]Edge, 0, len(tm)*len(tm))
	for i, row := range tm {
		for j, p := range row {
			if p == 0 {
				continue
			}

			edges = append(edges, Edge{
				ID:     fmt.Sprintf("edge-%d-%d", i, j),
				Label:  fmt.Sprintf("%.6f", p),
				Size:   3,
				Source: fmt.Sprintf("node-%d", i),
				Target: fmt.Sprintf("node-%d", j),
				Color:  "#ccc",
				Type:   "curvedArrow",
			})
		}
	}

	return Graph{
		Nodes: nodes,
		Edges: edges,
	}
}
