package stochastic

import (
	"math"
	"testing"
)

func BenchmarkBuildImplementation(b *testing.B) {
	m := 0.5
	n := 12
	h := 0.025

	mFn := func(h float64, m float64) func(int) float64 {
		return func(i int) float64 {
			return m
		}
	}(h, m)

	kFn := func(h float64) func(int, int) float64 {
		return func(i int, i2 int) float64 {
			return 1 / (1 + math.Pow(GetT(i, h)-GetT(i2, h), 2))
		}
	}(h)

	for i := 0; i < b.N; i += 1 {
		_ = BuildImplementation(mFn, kFn, n, true)
	}
}

func BenchmarkBuildImplementationGenerator(b *testing.B) {
	m := 0.5
	n := 12
	h := 0.025

	mFn := func(h float64, m float64) func(int) float64 {
		return func(i int) float64 {
			return m
		}
	}(h, m)

	kFn := func(h float64) func(int, int) float64 {
		return func(i int, i2 int) float64 {
			return 1 / (1 + math.Pow(GetT(i, h)-GetT(i2, h), 2))
		}
	}(h)

	implementationGenerator := BuildImplementationGenerator(mFn, kFn, n, true)

	for i := 0; i < b.N; i += 1 {
		_ = implementationGenerator()
	}
}
