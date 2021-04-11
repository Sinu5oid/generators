package chain

import (
	"gonum.org/v1/gonum/floats"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	g := NewGenerator([]float64{0, 0.3, 0, 0.6, 0.1, 0})

	expected := []segment{
		{lb: 0, rb: 0.3, v: 1},
		{lb: 0.3, rb: 0.9, v: 3},
		{lb: 0.9, rb: 1, v: 4},
	}

	for i, c := range expected {
		if !floats.EqualApprox([]float64{g.segments[i].lb}, []float64{c.lb}, tol) {
			t.Errorf("lb: expected %f got %f", g.segments[i].lb, c.lb)
			continue
		}

		if !floats.EqualApprox([]float64{g.segments[i].rb}, []float64{c.rb}, tol) {
			t.Errorf("rb: expected %f got %f", g.segments[i].rb, c.rb)
			continue
		}

		if g.segments[i].v != c.v {
			t.Errorf("v: expected %d got %d", g.segments[i].v, c.v)
			continue
		}
	}
}
