package rmd

import (
	"github.com/aviddiviner/go-funcache"
	"math"
	"math/rand"
)

type Engine struct {
	n            int
	sigma        float64
	sigmaSquared float64
	h            float64
	cache        *funcache.Cache
}

func NewEngine(n int, sigma, h float64) *Engine {
	return &Engine{
		n:            n,
		sigma:        sigma,
		sigmaSquared: sigma * sigma,
		h:            h,
		cache:        funcache.NewInMemCache(),
	}
}

func (e *Engine) NextImpl() []float64 {
	left := float64(0)
	right := rand.NormFloat64() * e.sigma

	buffer := make([]float64, 0)
	buffer = append(buffer, left)
	buffer = append(buffer, e.x(left, right, 1)...)
	buffer = append(buffer, right)
	return buffer
}

func (e *Engine) x(left, right float64, n int) []float64 {
	// recursion reached
	if n > e.n {
		return make([]float64, 0, 0)
	}

	middle := (left+right)/2 + rand.NormFloat64()*e.disp(n)
	leftPart := e.x(left, middle, n+1)
	rightPart := e.x(middle, right, n+1)

	buffer := make([]float64, 0)
	buffer = append(buffer, leftPart...)
	buffer = append(buffer, middle)
	buffer = append(buffer, rightPart...)

	return buffer
}

func (e *Engine) disp(n int) float64 {
	if n == 0 {
		return e.sigma
	}

	return e.cache.Cache(n, func() interface{} {
		return math.Sqrt((e.sigmaSquared * (1 - math.Pow(2, 2*e.h-2))) / math.Pow(2*float64(n), 2*e.h))
	}).(float64)
}

func (e *Engine) MeanE(impls [][]float64) []float64 {
	res := make([]float64, 0, len(impls[0]))
	for k := 0; k < len(impls[0]); k++ {
		meanAt := float64(0)
		for i := 0; i < len(impls); i++ {
			meanAt += impls[i][k]
		}
		meanAt /= float64(len(impls))

		res = append(res, meanAt)
	}

	return res
}

func (e *Engine) MeanT(impls [][]float64) []float64 {
	res := make([]float64, len(impls[0]), len(impls[0]))

	return res
}

func (e *Engine) DispE(impls [][]float64, meansE []float64) []float64 {
	res := make([]float64, 0, len(impls[0]))

	for k := 0; k < len(impls[0]); k++ {
		dispAt := float64(0)

		for i := 0; i < len(impls); i++ {
			dispAt += math.Pow(impls[i][k]-meansE[k], 2)
		}

		dispAt /= float64(len(impls) - 1)

		res = append(res, dispAt)
	}

	return res
}

func (e *Engine) DispT(impls [][]float64) []float64 {
	res := make([]float64, 0, len(impls[0]))

	for k := 0; k < len(impls[0]); k++ {
		dispAt := e.sigmaSquared * math.Pow(float64(k)/math.Pow(2, float64(e.n)), 2*e.h)
		res = append(res, dispAt)
	}

	return res
}

func (e *Engine) CorellE(impls [][]float64, dispsE []float64) [][]float64 {
	res := make([][]float64, 0, len(impls[0]))

	for k := 1; k < len(impls[0])-1; k++ {
		buf := make([]float64, 0, len(impls[0]))
		for j := 1; j < len(impls[0])-1-k; j++ {
			corellAt := float64(0)
			for i := 0; i < len(impls); i++ {
				impl := impls[i]
				interm := (impl[k+1] - impl[k]) * (impl[k+j+1] - impl[k+j])
				corellAt += interm
			}

			//corellAt /= -dispsE[1] * float64(len(impls)-1)
			//buf = append(buf, corellAt + 0.1)
			corellAt /= dispsE[1] * float64(len(impls)-1)
			buf = append(buf, corellAt)
		}
		res = append(res, padLeft(buf, len(impls[0])))
	}

	return res
}

func (e *Engine) CorellT(impls [][]float64, H float64) [][]float64 {
	res := make([][]float64, 0, len(impls[0]))

	for k := 1; k < len(impls[0])-1; k++ {
		buf := make([]float64, 0, len(impls[0]))
		for j := 1; j < len(impls[0])-1-k; j++ {
			buf = append(buf, (math.Pow(float64(j+1), 2*H)-2*math.Pow(float64(j), 2*H)+math.Pow(float64(j-1), 2*H))/2)
		}

		res = append(res, padLeft(buf, len(impls[0])))
	}

	return res
}

func padLeft(src []float64, targetLength int) []float64 {
	res := make([]float64, targetLength, targetLength)

	delta := targetLength - len(src)
	for i := 0; i < len(src); i++ {
		res[i+delta] = src[i]
	}

	return res
}
