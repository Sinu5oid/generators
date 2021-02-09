package stochastic

import (
	"fmt"
	"math"
	"math/rand"
)

type correlationFn func(int, int) float64
type meanFn func(int) float64

// buildModel
//
// m - mean function m(t)
//
// k - correlation function K(t, t')
//
// h - step
//
// n - time steps count
func buildModel(m meanFn, k correlationFn, n int) {
	devs := make([]float64, 0, n)
	funcs := make([][]float64, 0, n)

	for i := 0; i < n; i += 1 {
		devs = append(devs, *getDev(k, i, &devs, &funcs))
		funcs = append(funcs, *getFuncRow(k, i, n, &devs, &funcs))
	}

	randoms := make([]float64, 0, n)
	for i := 0; i < n; i++ {
		randoms = append(randoms, rand.NormFloat64()*math.Sqrt(devs[i]))
	}

	//randoms := []float64{-0.234158, -0.428654, 0.50671, 0.257188, -0.602093}
	impl := *getImpl(m, n, &funcs, &randoms)

	fmt.Println("Devs:")
	fmt.Println(devs)
	fmt.Println("Funcs:")
	fmt.Println(funcs)
	fmt.Println("Impls:")
	fmt.Println(impl)
}

func getDev(k correlationFn, i int, prevDevs *[]float64, prevFuncs *[][]float64) *float64 {
	res := k(i, i)

	for k := 0; k <= i-1; k += 1 {
		res = res - math.Pow((*prevFuncs)[k][i], 2)*(*prevDevs)[k]
	}

	return &res
}

func getFuncRow(k correlationFn, i int, n int, prevDevs *[]float64, prevFuncs *[][]float64) *[]float64 {
	funcs := make([]float64, 0, n)
	for j := 0; j < n; j += 1 {
		funcs = append(funcs, getFunc(k, i, j, prevDevs, prevFuncs))
	}

	return &funcs
}

func getFunc(k correlationFn, i int, j int, prevDevs *[]float64, prevFuncs *[][]float64) float64 {
	switch {
	case i == j:
		return 1
	case i > j:
		return 0
	default:
		dividend := k(i, j)
		for x := 0; x < i; x += 1 {
			dividend = dividend - (*prevFuncs)[x][i]*(*prevFuncs)[x][j]*(*prevDevs)[x]
		}

		return dividend / (*prevDevs)[i]
	}
}

func getImpl(m meanFn, n int, funcs *[][]float64, randoms *[]float64) *[]float64 {
	impls := make([]float64, 0, n)

	for i := 0; i < n; i += 1 {
		impl := m(i)

		for k := 0; k <= i; k += 1 {
			impl += (*randoms)[k] * (*funcs)[k][i]
		}

		impls = append(impls, impl)
	}

	return &impls
}

func getT(i int, h float64) float64 {
	return float64(i) * h
}

func BuildModel() {
	m := 1.0
	n := 5
	h := 0.25

	mFn := func(h float64) func(int) float64 {
		return func(i int) float64 {
			return m
		}
	}(h)

	kFn := func(h float64) func(int, int) float64 {
		return func(i int, i2 int) float64 {
			return 1 / (1 + 5*math.Abs(getT(i, h)-getT(i2, h)))
		}
	}(h)

	buildModel(mFn, kFn, n)
}
