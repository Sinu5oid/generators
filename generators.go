package generators

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
)

type IntGenerator interface{ Int() int }

type Float64Generator interface{ Float64() float64 }

type ExpFloat64Generator interface{ ExpFloat64() float64 }

type NormFloat64Generator interface{ NormFloat64() float64 }

type CongruentialGenerator struct {
	n       int
	m       int
	a       int
	initial int
	current int
}

func NewCongruentialGenerator(modulus int, multiplier int, additiveComponent int, initialValue int) *CongruentialGenerator {
	return &CongruentialGenerator{
		n:       modulus,
		m:       multiplier,
		a:       additiveComponent,
		current: initialValue,
		initial: initialValue,
	}
}

func (cg *CongruentialGenerator) String() string {
	d := make(map[string]interface{}, 4)

	d["distributionName"] = "congruential"
	d["modulus"] = cg.n
	d["multiplier"] = cg.m
	d["additiveComponent"] = cg.a
	d["initialValue"] = cg.initial

	if b, err := json.Marshal(d); err != nil {
		return fmt.Sprintf("%v", d)
	} else {
		return string(b)
	}
}

func (cg *CongruentialGenerator) Int() int {
	val := (cg.current*cg.m + cg.a) % cg.n
	cg.current = val

	return val
}

type UniformGenerator struct {
	g IntGenerator
	m int
}

func (ug *UniformGenerator) String() string {
	d := make(map[string]interface{}, 2)

	d["distributionName"] = "uniform"
	d["modulus"] = ug.m

	if b, err := json.Marshal(d); err != nil {
		return fmt.Sprintf("%v", d)
	} else {
		return string(b)
	}
}

func NewUniformGeneratorDefault() *UniformGenerator {
	return &UniformGenerator{g: rand.New(rand.NewSource(rand.Int63())), m: math.MaxInt32}
}

func NewUniformGenerator(generator IntGenerator, modulus int) *UniformGenerator {
	return &UniformGenerator{
		g: generator,
		m: modulus,
	}
}

func (ug *UniformGenerator) Float64() float64 {
	if ug.m == 0 {
		panic("modulus is 0")
	}

	return float64(ug.g.Int()) / float64(ug.m)
}

type ExponentialGenerator struct {
	g Float64Generator
	l float64
}

func (eg *ExponentialGenerator) String() string {
	d := make(map[string]interface{}, 2)

	d["distributionName"] = "exponential"
	d["rate"] = eg.l

	if b, err := json.Marshal(d); err != nil {
		return fmt.Sprintf("%v", d)
	} else {
		return string(b)
	}
}

func NewExponentialGeneratorDefault() *ExponentialGenerator {
	return &ExponentialGenerator{g: rand.New(rand.NewSource(rand.Int63())), l: 1}
}

func NewExponentialGenerator(generator Float64Generator, rate float64) *ExponentialGenerator {
	return &ExponentialGenerator{g: generator, l: rate}
}

func (eg *ExponentialGenerator) ExpFloat64() float64 {
	return -eg.l * math.Log(eg.g.Float64())
}

type NormalGenerator struct {
	g      Float64Generator
	stdDev float64
	mean   float64
}

func (ng *NormalGenerator) String() string {
	d := make(map[string]interface{}, 3)

	d["distributionName"] = "normal"
	d["standardDeviation"] = ng.stdDev
	d["mean"] = ng.mean

	if b, err := json.Marshal(d); err != nil {
		return fmt.Sprintf("%v", d)
	} else {
		return string(b)
	}
}

func NewNormalGeneratorDefault() *NormalGenerator {
	return &NormalGenerator{g: rand.New(rand.NewSource(rand.Int63())), stdDev: 1, mean: 0}
}

func NewNormalGenerator(generator Float64Generator, standardDeviation float64, mean float64) *NormalGenerator {
	return &NormalGenerator{g: generator, stdDev: standardDeviation, mean: mean}
}

func (ng *NormalGenerator) NormFloat64() float64 {
	v1 := 2*ng.g.Float64() - 1
	v2 := 2*ng.g.Float64() - 1
	S := math.Pow(v1, 2) + math.Pow(v2, 2)

	if S >= 1 {
		return ng.NormFloat64()
	}

	intermediate := math.Sqrt(-2/S*math.Log(S)) * ng.stdDev

	return intermediate*v1 + ng.mean
}
