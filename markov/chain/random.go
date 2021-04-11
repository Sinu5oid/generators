package chain

import (
	"golang.org/x/exp/rand"
)

type segment struct {
	lb float64
	rb float64
	v  int
}

type Generator struct {
	segments []segment
}

func NewGenerator(probs []float64) *Generator {
	segments := make([]segment, 0, len(probs))

	lb := float64(0)
	rb := float64(0)

	for i, p := range probs {
		if p == 0 {
			continue
		}

		lb = rb
		rb = lb + p

		segments = append(segments, segment{
			lb: lb,
			rb: rb,
			v:  i,
		})
	}

	return &Generator{segments: segments}
}

func (g Generator) Next() int {
	rnd := rand.Float64()

	for _, s := range g.segments {
		if s.lb <= rnd && s.rb > rnd {
			return s.v
		}
	}

	return g.segments[len(g.segments)-1].v
}
