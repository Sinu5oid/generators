package stochastic

import (
	"math"
	"math/rand"
)

type CorrelationFn func(int, int) float64
type MeanFn func(int) float64

// buildImplementation
//
// m - mean function m(t)
//
// k - correlation function K(t, t')
//
// n - time steps count
func buildImplementation(m MeanFn, k CorrelationFn, n int, trySafeMath bool) *[]float64 {
	devs := make([]float64, 0, n)
	funcs := make([][]float64, 0, n)

	for i := 0; i < n; i += 1 {
		devs = append(devs, *getDev(k, i, &devs, &funcs, trySafeMath))
		funcs = append(funcs, *getFuncRow(k, i, n, &devs, &funcs, trySafeMath))
	}

	randoms := make([]float64, 0, n)
	for i := 0; i < n; i += 1 {
		randoms = append(randoms, rand.NormFloat64()*math.Sqrt((devs)[i]))
	}

	return getImpl(m, n, &funcs, &randoms)
}

// buildImplementationTemplate
//
// k - correlation function K(t, t')
//
// n - time steps count
func buildImplementationTemplate(
	k CorrelationFn,
	n int,
	trySafeMath bool,
) (*[]float64, *[][]float64) {
	devs := make([]float64, 0, n)
	funcs := make([][]float64, 0, n)

	for i := 0; i < n; i += 1 {
		devs = append(devs, *getDev(k, i, &devs, &funcs, trySafeMath))
		funcs = append(funcs, *getFuncRow(k, i, n, &devs, &funcs, trySafeMath))
	}

	return &devs, &funcs
}

func getDev(
	k CorrelationFn,
	i int,
	prevDevs *[]float64,
	prevFuncs *[][]float64,
	trySafeMath bool,
) *float64 {
	kVal := k(i, i)
	sum := 0.0

	for x := 0; x < i; x += 1 {
		sum += math.Pow((*prevFuncs)[x][i], 2) * (*prevDevs)[x]
	}

	res := kVal - sum

	// underflow detection & protection
	if (sum > kVal || math.IsNaN(res)) && trySafeMath {
		res = 0
	}

	return &res
}

func getFuncRow(
	k CorrelationFn,
	i, n int,
	prevDevs *[]float64,
	prevFuncs *[][]float64,
	trySafeMath bool,
) *[]float64 {
	funcs := make([]float64, 0, n)
	for j := 0; j < n; j += 1 {
		funcs = append(funcs, getFunc(k, i, j, prevDevs, prevFuncs, trySafeMath))
	}

	return &funcs
}

func getFunc(
	k CorrelationFn,
	i, j int,
	prevDevs *[]float64,
	prevFuncs *[][]float64,
	trySafeMath bool,
) float64 {
	switch {
	case i == j:
		return 1
	case i > j:
		return 0
	default:
		if trySafeMath && (*prevDevs)[i] == 0 {
			return 0
		}

		dividend := k(i, j)
		for x := 0; x < i; x += 1 {
			dividend = dividend - (*prevFuncs)[x][i]*(*prevFuncs)[x][j]*(*prevDevs)[x]
		}

		return dividend / (*prevDevs)[i]
	}
}

func getImpl(
	m MeanFn,
	n int,
	funcs *[][]float64,
	randoms *[]float64,
) *[]float64 {
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

func BuildImplementationGenerator(
	mFn MeanFn,
	kFn CorrelationFn,
	n int,
	trySafeMath bool,
) func() *[]float64 {
	devs, funcs := buildImplementationTemplate(kFn, n, trySafeMath)

	return func() *[]float64 {
		randoms := make([]float64, 0, n)
		for i := 0; i < n; i += 1 {
			randoms = append(randoms, rand.NormFloat64()*math.Sqrt((*devs)[i]))
		}

		return getImpl(mFn, n, funcs, &randoms)
	}
}

func BuildImplementation(mFn MeanFn, kFn CorrelationFn, n int, trySafeMath bool) *[]float64 {
	return buildImplementation(mFn, kFn, n, trySafeMath)
}

func GetT(i int, h float64) float64 {
	return float64(i) * h
}
