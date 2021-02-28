package main

import (
	"fmt"
	"github.com/Sinu5oid/generators/stochastic"
	"sort"
	"sync"
	"time"
)

func run(
	n int,
	h float64,
	N int,
	mFn stochastic.MeanFn,
	kFn stochastic.CorrelationFn,
	implsToDisplay int,
	outputFolder string,
	numRunners int,
	trySafeMath bool,
	useSingleTemplate bool,
) {
	fmt.Println("allowed to use", numRunners, "routines")

	implsInChan := make(chan struct{}, N)
	implsChan := make(chan *[]float64, N)

	started := time.Now()

	fmt.Println("[implementation gens] started", numRunners, "implementation generators")
	var wg sync.WaitGroup

	wg.Add(numRunners)
	if useSingleTemplate {
		implementationGenerator := stochastic.BuildImplementationGenerator(mFn, kFn, n, trySafeMath)

		for i := 0; i < numRunners; i += 1 {
			go func(inCh chan struct{}, outCh chan *[]float64, wg *sync.WaitGroup) {
				defer wg.Done()

				for range inCh {
					outCh <- implementationGenerator()
				}
			}(implsInChan, implsChan, &wg)
		}
	} else {
		for i := 0; i < numRunners; i += 1 {
			go func(inCh chan struct{}, outCh chan *[]float64, wg *sync.WaitGroup) {
				defer wg.Done()

				for range inCh {
					outCh <- stochastic.BuildImplementation(mFn, kFn, n, trySafeMath)
				}
			}(implsInChan, implsChan, &wg)
		}
	}

	// discharge channels
	for i := 0; i < N; i += 1 {
		implsInChan <- struct{}{}
	}

	close(implsInChan)
	wg.Wait()
	close(implsChan)

	impls := make([][]float64, 0, N)
	for impl := range implsChan {
		impls = append(impls, *impl)
	}

	fmt.Println("[implementation gens] finished", numRunners, "implementation generators in", time.Since(started))

	meansIdealChan := make(chan *float64, n)
	funcsIdealChan := make(chan *[]*float64, n)
	meansObservedChan := make(chan *float64, n)
	funcsObservedChan := make(chan *[]*float64, n)

	started = time.Now()
	fmt.Println("[stat analysis] started")
	wg.Add(4)
	go func(ch chan *float64, wg *sync.WaitGroup, mFn stochastic.MeanFn, n int) {
		defer wg.Done()
		fmt.Println("\t[means ideal] started")
		for i := 0; i < n; i += 1 {
			res := mFn(i)
			ch <- &res
		}
		fmt.Println("\t[means ideal] finished")
		close(ch)
	}(meansIdealChan, &wg, mFn, n)

	go func(ch chan *[]*float64, wg *sync.WaitGroup, kFn stochastic.CorrelationFn, n int) {
		defer wg.Done()
		fmt.Println("\t[funcs ideal] started")
		for i := 0; i < n; i += 1 {
			ch <- getFuncsIdealRow(kFn, n, i)
		}
		fmt.Println("\t[funcs ideal] finished")
		close(ch)
	}(funcsIdealChan, &wg, kFn, n)

	go func(ch chan *float64, wg *sync.WaitGroup, impls *[][]float64, n int, N int) {
		defer wg.Done()
		fmt.Println("\t[means observed] started")
		for i := 0; i < n; i += 1 {
			ch <- getMeanObserved(i, N, impls)

		}
		fmt.Println("\t[means observed] finished")
		close(ch)
	}(meansObservedChan, &wg, &impls, n, N)

	go func(ch chan *[]*float64, wg *sync.WaitGroup, impls *[][]float64, n int, N int) {
		defer wg.Done()
		fmt.Println("\t[funcs observed] started")
		for i := 0; i < n; i += 1 {
			ch <- getFuncsObservedRow(N, n, i, impls)
		}
		fmt.Println("\t[funcs observed] finished")
		close(ch)
	}(funcsObservedChan, &wg, &impls, n, N)

	wg.Wait()

	meansIdeal := make([]*float64, 0, n)
	meansObserved := make([]*float64, 0, n)
	funcsObserved := make([]*[]*float64, 0, n)
	funcsIdeal := make([]*[]*float64, 0, n)

	for i := 0; i < n; i += 1 {
		meansIdeal = append(meansIdeal, <-meansIdealChan)
		meansObserved = append(meansObserved, <-meansObservedChan)
		funcsObserved = append(funcsObserved, <-funcsObservedChan)
		funcsIdeal = append(funcsIdeal, <-funcsIdealChan)
	}

	fmt.Println("[stat analysis] finished in", time.Since(started))

	taskChan := make(chan GenericTask, n)

	started = time.Now()
	wg.Add(numRunners)
	fmt.Println("[plotter tasks] started", numRunners, "task consumers")
	for i := 0; i < numRunners; i += 1 {
		go func(i int, wg *sync.WaitGroup, ch chan GenericTask) {
			defer wg.Done()

			for task := range ch {
				if err := task.BuildPlot(); err != nil {
					continue
				}
			}
		}(i, &wg, taskChan)
	}

	meansInterm := make([]float64, 0, len(meansIdeal))
	for i := 0; i < len(meansIdeal); i += 1 {
		meansInterm = append(meansInterm, *meansIdeal[i])
	}

	sort.Float64s(meansInterm)

	meanYMin := meansInterm[0] - meansInterm[0]*0.1
	meanYMax := meansInterm[len(meansInterm)-1] + meansInterm[len(meansInterm)-1]*0.1
	taskChan <- CreateTask(
		h,
		"mean",
		0,
		"t",
		"m(t)",
		&meansIdeal,
		&meansObserved,
		&meanYMin,
		&meanYMax,
		outputFolder,
	)

	sliceToDisplay := impls[0:implsToDisplay]
	taskChan <- CreateSimpleTask(
		h,
		"implementations",
		"t",
		"xi(t)",
		&sliceToDisplay,
		outputFolder,
	)
	for i := 0; i < n; i += 1 {
		taskChan <- CreateTask(
			h,
			"correlation",
			i,
			"t'",
			fmt.Sprintf("K(t%d, t')", i),
			funcsIdeal[i],
			funcsObserved[i],
			nil,
			nil,
			outputFolder,
		)
	}

	close(taskChan)
	wg.Wait()
	fmt.Println("[plotter tasks] finished", numRunners, "task consumers")
}
