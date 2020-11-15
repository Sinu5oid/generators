package main

import (
	"fmt"
	"github.com/Sinu5oid/generators"
	"github.com/Sinu5oid/generators/stat"
	"math"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	maxIterations := int(math.Pow(10, 7))

	rate := 1.0
	stdDev := 1.0
	mean := 2.0

	intervalsCount := 100
	confidenceLevel := 0.05

	modulus1 := int(math.Pow(2, 32))
	cg := generators.NewCongruentialGenerator(modulus1, 1103515245, 12345, 0)
	ug := generators.NewUniformGenerator(cg, modulus1)
	eg := generators.NewExponentialGenerator(ug, rate)

	modulus2 := int(math.Pow(2, 32))
	cg2 := generators.NewCongruentialGenerator(modulus2, 134775813, 1, 3)
	ug2 := generators.NewUniformGenerator(cg2, modulus2)
	ng := generators.NewNormalGenerator(ug, ug2, stdDev, mean)

	modulus3 := int(math.Pow(2, 31)) - 1
	cg3 := generators.NewCongruentialGenerator(modulus3, 2147483629, 2147483587, 255)
	ug3 := generators.NewUniformGenerator(cg3, modulus3)

	e := make(chan []float64, 1)
	n := make(chan []float64, 1)
	u := make(chan []float64, 1)

	go getDistribution(eg.Name(), eg.String(), eg.ExpFloat64, maxIterations, &wg, e)
	go getDistribution(ng.Name(), ng.String(), ng.NormFloat64, maxIterations, &wg, n)
	go getDistribution(ug3.Name(), ug3.String(), ug3.Float64, maxIterations, &wg, u)

	wg.Wait()

	nd := <-n
	ed := <-e
	ud := <-u

	runNormalDistributionAnalysis(nd, stdDev, mean, intervalsCount, confidenceLevel)
	runExponentialDistributionAnalysis(ed, rate, intervalsCount, confidenceLevel)
	runUniformDistributionAnalysis(ud, intervalsCount, confidenceLevel)

	wg.Wait()
}

func runNormalDistributionAnalysis(distribution []float64, stdDev float64, mean float64, intervalsCount int, confidenceLevel float64) {
	fmt.Println("###normal distribution test started")
	analysis := stat.NewStatisticAnalysis(distribution, intervalsCount, confidenceLevel)

	chiObserved, chiCritical := analysis.TestPearsonNormal(stdDev, mean)
	if chiObserved > chiCritical {
		fmt.Println("[NRM] distribution is not of normal type because", chiObserved, ">", chiCritical)
	} else {
		fmt.Println("[NRM] distribution Pearson test passed", chiObserved, "<=", chiCritical)
	}
	fmt.Println("###normal distribution test finished")
}

func runExponentialDistributionAnalysis(distribution []float64, rate float64, intervalsCount int, confidenceLevel float64) {
	fmt.Println("###exponential distribution test started")
	analysis := stat.NewStatisticAnalysis(distribution, intervalsCount, confidenceLevel)

	chiObserved, chiCritical := analysis.TestPearsonExp(rate)
	if chiObserved > chiCritical {
		fmt.Println("[EXP] distribution is not of exponential type because", chiObserved, ">", chiCritical)
	} else {
		fmt.Println("[EXP] distribution Pearson test passed", chiObserved, "<=", chiCritical)
	}
	fmt.Println("###exponential distribution test finished")
}

func runUniformDistributionAnalysis(distribution []float64, intervalsCount int, confidenceLevel float64) {
	fmt.Println("###uniform distribution test started")
	analysis := stat.NewStatisticAnalysis(distribution, intervalsCount, confidenceLevel)

	chiObserved, chiCritical := analysis.TestPearsonUniform()
	if chiObserved > chiCritical {
		fmt.Println("[UNI] distribution is not of uniform type because", chiObserved, ">", chiCritical)
	} else {
		fmt.Println("[UNI] distribution Pearson test passed", chiObserved, "<=", chiCritical)
	}
	fmt.Println("###uniform distribution test finished")
}

func getDistribution(distributionName, characteristics string, source func() float64, maxIterations int, wg *sync.WaitGroup, result chan []float64) {
	wg.Add(1)

	defer func() {
		fmt.Printf("%q generator done\n", distributionName)
		wg.Done()
	}()

	fmt.Printf("Running %q generator, target values count: %d\n", distributionName, maxIterations)
	fmt.Printf("Characteristics: %s\n", characteristics)
	generatedValues := make([]float64, 0, maxIterations)

	for i := 0; i < maxIterations; i += 1 {
		generatedValues = append(generatedValues, source())
	}

	result <- generatedValues
}
