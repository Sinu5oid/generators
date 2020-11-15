package stat

import (
	"fmt"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
	"sort"
	"time"
)

type StatisticAnalysis struct {
	// distribution result
	source []float64

	// intervals
	intervals []interval

	// confidence level
	alpha float64

	max float64
	min float64
}

func NewStatisticAnalysis(distribution []float64, intervalsCount int, confidenceLevel float64) StatisticAnalysis {
	if len(distribution) < 1 {
		panic("distribution length is less than 1")
	}

	sort.Float64s(distribution)
	min := distribution[0]
	max := distribution[len(distribution)-1]

	intervalsDelta := (max - min) / float64(intervalsCount)

	intervalsBounds := make([]float64, 0)
	for i := min; i <= max+0.01; i += intervalsDelta {
		intervalsBounds = append(intervalsBounds, i)
	}

	intervals := merge(split(distribution, intervalsBounds))

	totalValues := 0

	for i, interval := range intervals {
		totalValues += len(interval.values)
		fmt.Printf("interval %d\t\t[%+.6f, %+.6f)\t\t%d values\n", i, interval.leftBound, interval.rightBound, len(interval.values))
	}

	fmt.Printf("total values: %d of %d\n", totalValues, len(distribution))

	return StatisticAnalysis{
		source:    distribution,
		alpha:     confidenceLevel,
		intervals: intervals,
		min:       min,
		max:       max,
	}
}

func (s StatisticAnalysis) TestPearsonNormal(sigma float64, alpha float64) (float64, float64) {
	probabilities := make([]float64, 0, len(s.intervals))

	if len(s.intervals) < 1 {
		panic("intervals length is less than 1")
	}

	for i := 1; i < len(s.intervals); i += 1 {
		p := normalCDF(s.intervals[i].leftBound, sigma, alpha) -
			normalCDF(s.intervals[i-1].leftBound, sigma, alpha)

		fmt.Printf(
			"p(%d) = %.6f\t\t\t[%+.6f, %+.6f)\n",
			i-1,
			p,
			s.intervals[i-1].leftBound,
			s.intervals[i].leftBound,
		)

		probabilities = append(probabilities, p)
	}

	lastP := normalCDF(s.intervals[len(s.intervals)-1].rightBound, sigma, alpha) -
		normalCDF(s.intervals[len(s.intervals)-1].leftBound, sigma, alpha)
	probabilities = append(
		probabilities,
		lastP,
	)

	fmt.Printf(
		"p(%d) = %.6f\t\t[%+.6f, %+.6f)\n",
		len(s.intervals)-1,
		lastP,
		s.intervals[len(s.intervals)-1].leftBound,
		s.intervals[len(s.intervals)-1].rightBound,
	)

	chiObserved := 0.0
	cumulativeProbability := 0.0

	for i, p := range probabilities {
		cumulativeProbability += p
		chiObserved += math.Pow(
			float64(len(s.intervals[i].values))/float64(len(s.source))-p,
			2,
		) / p
	}

	fmt.Println("sum(p) = ", cumulativeProbability)

	chiObserved = chiObserved * float64(len(s.source))

	chiCritical := distuv.ChiSquared{
		K:   float64(len(s.intervals) - 1),
		Src: rand.NewSource(uint64(time.Now().UnixNano())),
	}

	return chiObserved, chiCritical.Quantile(1 - s.alpha)
}

func (s StatisticAnalysis) TestPearsonExp(lambda float64) (float64, float64) {
	probabilities := make([]float64, 0, len(s.intervals))

	if len(s.intervals) < 1 {
		panic("intervals length is less than 1")
	}

	for i := 1; i < len(s.intervals); i += 1 {
		p := exponentialCDF(s.intervals[i].leftBound, lambda) -
			exponentialCDF(s.intervals[i-1].leftBound, lambda)

		fmt.Printf(
			"p(%d) = %.30f\t\t\t[%+.6f, %+.6f)\n",
			i-1,
			p,
			s.intervals[i-1].leftBound,
			s.intervals[i].leftBound,
		)

		probabilities = append(probabilities, p)
	}

	lastP := exponentialCDF(s.intervals[len(s.intervals)-1].rightBound, lambda) -
		exponentialCDF(s.intervals[len(s.intervals)-1].leftBound, lambda)
	probabilities = append(
		probabilities,
		lastP,
	)

	fmt.Printf(
		"p(%d) = %.30f\t\t[%+.6f, %+.6f)\n",
		len(s.intervals)-1,
		lastP,
		s.intervals[len(s.intervals)-1].leftBound,
		s.intervals[len(s.intervals)-1].rightBound,
	)

	chiObserved := 0.0
	cumulativeProbability := 0.0

	for i, p := range probabilities {
		cumulativeProbability += p
		chiObserved += math.Pow(
			float64(len(s.intervals[i].values))/float64(len(s.source))-p,
			2,
		) / p
	}

	fmt.Println("sum(p) = ", cumulativeProbability)

	chiObserved = chiObserved * float64(len(s.source))

	chiCritical := distuv.ChiSquared{
		K:   float64(len(s.intervals) - 1),
		Src: rand.NewSource(uint64(time.Now().UnixNano())),
	}

	return chiObserved, chiCritical.Quantile(1 - s.alpha)
}

func (s StatisticAnalysis) TestPearsonUniform() (float64, float64) {
	probabilities := make([]float64, 0, len(s.intervals))

	if len(s.intervals) < 1 {
		panic("intervals length is less than 1")
	}

	for i := 1; i < len(s.intervals); i += 1 {
		p := uniformCDF(s.intervals[i].leftBound, s.min, s.max) -
			uniformCDF(s.intervals[i-1].leftBound, s.min, s.max)

		fmt.Printf(
			"p(%d) = %.6f\t\t\t[%+.6f, %+.6f)\n",
			i-1,
			p,
			s.intervals[i-1].leftBound,
			s.intervals[i].leftBound,
		)

		probabilities = append(probabilities, p)
	}

	lastP := uniformCDF(s.intervals[len(s.intervals)-1].rightBound, s.min, s.max) -
		uniformCDF(s.intervals[len(s.intervals)-1].leftBound, s.min, s.max)
	probabilities = append(
		probabilities,
		lastP,
	)

	fmt.Printf(
		"p(%d) = %.6f\t\t[%+.6f, %+.6f)\n",
		len(s.intervals)-1,
		lastP,
		s.intervals[len(s.intervals)-1].leftBound,
		s.intervals[len(s.intervals)-1].rightBound,
	)

	chiObserved := 0.0
	cumulativeProbability := 0.0

	for i, p := range probabilities {
		cumulativeProbability += p
		chiObserved += math.Pow(
			float64(len(s.intervals[i].values))/float64(len(s.source))-p,
			2,
		) / p
	}

	fmt.Println("sum(p) = ", cumulativeProbability)

	chiObserved = chiObserved * float64(len(s.source))

	chiCritical := distuv.ChiSquared{
		K:   float64(len(s.intervals) - 1),
		Src: rand.NewSource(uint64(time.Now().UnixNano())),
	}

	return chiObserved, chiCritical.Quantile(1 - s.alpha)
}
