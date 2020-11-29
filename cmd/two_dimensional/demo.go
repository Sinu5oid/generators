package main

import (
	"fmt"
	"github.com/Sinu5oid/generators"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"math"
)

func main() {
	maxIterations := int(math.Pow(10, 4))

	modulus1 := int(math.Pow(2, 32))
	cg := generators.NewCongruentialGenerator(modulus1, 1103515245, 12345, 0)
	ug := generators.NewUniformGenerator(cg, modulus1)

	modulus2 := int(math.Pow(2, 32))
	cg2 := generators.NewCongruentialGenerator(modulus2, 134775813, 1, 3)
	ug2 := generators.NewUniformGenerator(cg2, modulus2)

	tdg := generators.NewTwoDimensionalGenerator(ug, ug2, 1, 1, 0, 0, 0.1)
	tdg2 := generators.NewTwoDimensionalGenerator(ug, ug2, 1, 1, 0, 0, 0.5)
	tdg3 := generators.NewTwoDimensionalGenerator(ug, ug2, 1, 1, 0, 0, 0.9)

	distr := make(generators.FloatPairs, 0, maxIterations)
	distr2 := make(generators.FloatPairs, 0, maxIterations)
	distr3 := make(generators.FloatPairs, 0, maxIterations)

	for i := 0; i < maxIterations; i += 1 {
		distr = append(distr, tdg.TwoDimensionalFloat64s())
		distr2 = append(distr2, tdg2.TwoDimensionalFloat64s())
		distr3 = append(distr3, tdg3.TwoDimensionalFloat64s())
	}

	p, err := plot.New()
	if err != nil {
		fmt.Printf("can't create plot: %s\n", err)
		return
	}
	p.Title.Text = "Distribution"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	s, err := plotter.NewScatter(distr)
	if err != nil {
		log.Panic(err)
	}
	s.GlyphStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	s.GlyphStyle.Radius = vg.Points(1)

	s2, err := plotter.NewScatter(distr2)
	if err != nil {
		log.Panic(err)
	}
	s2.GlyphStyle.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	s2.GlyphStyle.Radius = vg.Points(1)

	s3, err := plotter.NewScatter(distr3)
	if err != nil {
		log.Panic(err)
	}
	s3.GlyphStyle.Color = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	s3.GlyphStyle.Radius = vg.Points(1)

	p.Add(s, s2, s3)

	p.Legend.Add("r = 0.1", s)
	p.Legend.Add("r = 0.5", s2)
	p.Legend.Add("r = 0.9", s3)

	err = p.Save(10*vg.Inch, 10*vg.Inch, "scatter.png")
	if err != nil {
		log.Panic(err)
	}
}
