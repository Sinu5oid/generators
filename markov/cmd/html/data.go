package html

import (
	"github.com/Sinu5oid/generators/markov/chain"
	"github.com/Sinu5oid/generators/markov/cmd/diff"
)

type PageData struct {
	Graph           chain.Graph
	Implementations [][]int
	Diffs           [][]diff.Info
}
