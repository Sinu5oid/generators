package main

import (
	"fmt"
	"github.com/Sinu5oid/generators"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"math"
)

func main() {
	maxIterations := int(math.Pow(10, 5))
	modulus := int(math.Pow(2, 32))
	cg := generators.NewCongruentialGenerator(modulus, 1103515245, 12345, 0)
	ug := generators.NewUniformGenerator(cg, modulus)
	eg := generators.NewExponentialGenerator(ug, 345)
	ng := generators.NewNormalGenerator(ug, 1, 2)

	runDistributionAnalysis(cg,
		func(f func() int) func() float64 {
			return func() float64 {
				return float64(f())
			}
		}(cg.Int),
		maxIterations,
		25,
	)
	runDistributionAnalysis(ug, ug.Float64, maxIterations, 25)
	runDistributionAnalysis(eg, eg.ExpFloat64, maxIterations, 200)
	runDistributionAnalysis(ng, ng.NormFloat64, maxIterations, 75)
}

func runDistributionAnalysis(distribution generators.DistributionGenerator, source func() float64, maxIterations int, colCount int) {
	fmt.Printf("Running %q, target values count: %d\n", distribution.Name(), maxIterations)
	fmt.Printf("Characteristics: %s\n", distribution)
	generatedValues := make(plotter.Values, 0, maxIterations)

	for i := 0; i < maxIterations; i += 1 {
		generatedValues = append(generatedValues, source())
		i += 1
	}

	p, err := plot.New()
	if err != nil {
		fmt.Printf("can't create plot: %s\n", err)
		return
	}
	p.Title.Text = fmt.Sprintf("Histogram (%s)\n", distribution.Name())

	h, err := plotter.NewHist(generatedValues, colCount)
	if err != nil {
		fmt.Printf("can't create histogram: %s\n", err)
		return
	}
	h.Normalize(1)
	p.Add(h)

	filename := fmt.Sprintf("%s-hist.png", distribution.Name())

	if err := p.Save(10*vg.Inch, 5*vg.Inch, filename); err != nil {
		fmt.Printf("can't save file: %s\n", err)
		return
	}

	fmt.Printf("%q finished, artifact: %s\n\n", distribution.Name(), filename)
}
