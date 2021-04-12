package chain

import (
	"fmt"
	"github.com/aviddiviner/go-funcache"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/floats"
)

type Engine struct {
	tm    [][]float64
	s     int
	sc    int
	cache *funcache.Cache
}

func (e Engine) TransitionGraph() Graph {
	return e.cache.Wrap(func() interface{} {
		return NewGraph(e.tm)
	}).(Graph)
}

func (e *Engine) WithSteps(steps int) *Engine {
	e.sc = steps
	return e
}

func (e Engine) NextImpl() []int {
	res := make([]int, 0, e.sc+1)

	curr := e.s
	res = append(res, curr)
	for i := 0; i < e.sc; i++ {
		row := e.tm[curr]
		if floats.EqualApprox([]float64{row[curr]}, []float64{1}, 1e-25) {
			return res
		}

		curr = NewGenerator(row).Next()
		res = append(res, curr)
	}

	return res
}

func (e Engine) TProb(t int) []float64 {
	return e.cache.Cache(
		fmt.Sprintf("TProb-%d", t),
		func() interface{} {
			if t < 0 {
				res := make([]float64, len(e.tm), len(e.tm))
				res[e.s] = 1

				return res
			}

			return mustMultiplyMatrices([][]float64{e.TProb(t - 1)}, e.tm)[0]
		}).([]float64)
}

func NewEngine(tm [][]float64, s int) *Engine {
	if err := validateMatrix(tm); err != nil {
		panic(err)
	}

	if s > len(tm) {
		panic(errors.Wrap(ErrInvalidArguments, "s is out of bounds"))
	}

	var cache *funcache.Cache
	store, err := lru.New2Q(20)
	if err == nil {
		cache = funcache.New(store)
	} else {
		cache = funcache.NewInMemCache()
	}

	return &Engine{
		tm:    tm,
		s:     s,
		sc:    len(tm),
		cache: cache,
	}
}
