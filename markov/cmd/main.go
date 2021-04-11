package main

import (
	"fmt"
	"github.com/Sinu5oid/generators/markov/chain"
	"golang.org/x/exp/rand"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	rand.Seed(uint64(time.Now().UnixNano()))

	logger := log.New(os.Stdout, "", 0)
	started := time.Now()

	defer func() {
		logger.Println("finished in", time.Since(started))
	}()

	// transition matrix
	tm := [][]float64{
		{1, 0, 0, 0, 0},
		{7.0 / 17.0, 6.0 / 17.0, 0, 4.0 / 17.0, 0},
		{0, 0, 0, 10.0 / 13.0, 3.0 / 13.0},
		{0, 0, 1.0 / 2.0, 1.0 / 2.0, 0},
		{8.0 / 15.0, 0, 0, 7.0 / 15.0, 0},
	}

	// starting node
	s := 2

	// steps count
	sc := 25

	// implementations count
	ic := 1000000

	// cpus available
	cpus := runtime.NumCPU()

	e := chain.NewEngine(tm, s)

	// build transition graph (adjacency list)
	graph := e.TransitionGraph()
	fmt.Println("transition graph:", graph)

	// do [0-9] imitation x 1000 implementations
	e = e.WithSteps(sc)

	implsc := make(chan []int, ic)

	wg := sync.WaitGroup{}
	wg.Add(cpus)

	logger.Println("started implementations generation")
	// set up workers pool
	for runnerIndex := 0; runnerIndex < cpus; runnerIndex++ {
		its := ic / cpus
		if runnerIndex == cpus-1 {
			// compensate for extra iterations left
			its = (ic / cpus) + (ic % cpus)
		}

		go func(wg *sync.WaitGroup, e *chain.Engine, out chan<- []int, its int, index int) {
			defer wg.Done()
			defer func() {
				logger.Println("routine #", index, "finished")
			}()

			for i := 0; i < its; i++ {
				out <- e.NextImpl()
			}
		}(&wg, e, implsc, its, runnerIndex)
	}

	// discharge generated implementations
	impls := make([][]int, 0, ic)
	wg.Add(1)
	go func() {
		logger.Println("started discharging generated implementations")
		defer wg.Done()
		defer func() {
			logger.Println("finished discharging generated implementations")
		}()

		for i := 0; i < ic; i++ {
			impl, ok := <-implsc
			if !ok {
				break
			}
			impls = append(impls, impl)
		}
		close(implsc)
	}()

	// get theoretical p(t)
	logger.Println("started computing theoretical p(t)")
	tprobs := make([][]float64, 0, sc)
	for t := 0; t < sc; t++ {
		tprobs = append(tprobs, e.TProb(t))
	}
	logger.Println("finished computing theoretical p(t)")

	wg.Wait()

	logger.Println("implementations (up to first 20 elem. slice):")
	for i, impl := range impls[:20] {
		logger.Printf("#%2d %v Term.:%t\n", i, impl, len(impl) < sc)
	}
	logger.Println("-----")

	// get empiric p*(t)
	logger.Println("started computing empiric p*(t)")
	eprobs := make([][]float64, 0, sc)
	for t := 0; t < sc; t++ {
		eprobs = append(eprobs, chain.EProb(impls, len(tm), t, ic))
	}
	logger.Println("finished computing empiric p*(t)")

	logger.Println("comparison")
	for t := 0; t < sc; t++ {
		logger.Println("step #", t)
		for _, i := range diff(tprobs[t], eprobs[t]) {
			logger.Printf("t:\t%.6f\t|\te:\t%.6f\t(%+.6f)", i.t, i.e, i.diff)
		}
		logger.Println("-----")
	}
}
