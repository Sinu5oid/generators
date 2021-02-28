package main

import (
	"fmt"
	"github.com/Sinu5oid/generators/stochastic"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	globalStarted := time.Now()
	m := 0.5
	n := 12
	h := 0.025
	N := 1000000
	implsToDisplay := 6

	mFn := func(h float64, m float64) func(int) float64 {
		return func(i int) float64 {
			return m
		}
	}(h, m)

	kFn := func(h float64) func(int, int) float64 {
		return func(i int, i2 int) float64 {
			return 1 / (1 + math.Pow(stochastic.GetT(i, h)-stochastic.GetT(i2, h), 2))
		}
	}(h)

	var wg sync.WaitGroup
	wg.Add(1)
	cpusAvailable := runtime.NumCPU()
	impls := int(math.Min(float64(implsToDisplay), float64(N)))

	// using single template
	go startRunner(
		&wg,
		"stochastic/cmd/output-single-template",
		cpusAvailable,
		impls,
		n,
		h,
		N,
		mFn,
		kFn,
		true,
		true,
	)

	wg.Wait()

	fmt.Println("finished in", time.Since(globalStarted))
}
