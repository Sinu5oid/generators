package main

import (
	"github.com/Sinu5oid/generators/stochastic"
	"os"
	"sync"
)

func getMeanObserved(i int, N int, impls *[][]float64) *float64 {
	res := 0.0
	for k := 0; k < N; k += 1 {
		res += (*impls)[k][i]
	}

	res = res / float64(N)

	return &res
}

func getFuncObserved(n int, i int, j int, impls *[][]float64) *float64 {
	sij := 0.0
	si := 0.0
	sj := 0.0

	for k := 0; k < n; k += 1 {
		sij += (*impls)[k][i] * (*impls)[k][j]
		si += (*impls)[k][i]
		sj += (*impls)[k][j]
	}

	result := (sij - (si*sj)/float64(n)) / float64(n-1)
	return &result
}

func getFuncsObservedRow(N int, n int, i int, impls *[][]float64) *[]*float64 {
	funcsObserved := make([]*float64, 0, n)
	for j := 0; j < n; j += 1 {
		funcsObserved = append(funcsObserved, getFuncObserved(N, i, j, impls))
	}

	return &funcsObserved
}

func getFuncsIdealRow(kFn stochastic.CorrelationFn, n int, i int) *[]*float64 {
	funcsReal := make([]*float64, 0, n)
	for j := 0; j < n; j += 1 {
		res := kFn(i, j)
		funcsReal = append(funcsReal, &res)
	}

	return &funcsReal
}

func ensureFolderCreated(path string) {
	_ = os.RemoveAll(path)

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic("failed to create output folder")
	}
}

func startRunner(
	wg *sync.WaitGroup,
	outputFolder string,
	cpusAvailable, implsCount, n int,
	h float64,
	N int,
	mFn stochastic.MeanFn,
	kFn stochastic.CorrelationFn,
	trySafeMath, useSingleTemplate bool,
) {
	defer wg.Done()
	ensureFolderCreated(outputFolder)

	run(
		n,
		h,
		N,
		mFn,
		kFn,
		implsCount,
		outputFolder,
		cpusAvailable,
		trySafeMath,
		useSingleTemplate,
	)
}
