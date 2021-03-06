package generators

import (
	"encoding/json"
	"fmt"
	"gonum.org/v1/gonum/floats"
	"math"
	"math/rand"
)

type GeneratorName string

const (
	Congruential   GeneratorName = "congruential"
	Uniform        GeneratorName = "uniform"
	Exponential    GeneratorName = "exponential"
	Normal         GeneratorName = "normal"
	TwoDimensional GeneratorName = "two-dimensional"
)

type IntGenerator interface{ Int() int }

type Float64Generator interface{ Float64() float64 }

type ExpFloat64Generator interface{ ExpFloat64() float64 }

type NormFloat64Generator interface{ NormFloat64() float64 }

type TwoDimensionalFloat64Generator interface {
	TwoDimensionalFloat64s() (float64, float64)
}

type DistributionGenerator interface {
	String() string
	Name() string
}

type CongruentialGenerator struct {
	name    GeneratorName
	n       int
	m       int
	a       int
	initial int
	current int
}

func NewCongruentialGenerator(modulus int, multiplier int, additiveComponent int, initialValue int) *CongruentialGenerator {
	return &CongruentialGenerator{
		name:    Congruential,
		n:       modulus,
		m:       multiplier,
		a:       additiveComponent,
		current: initialValue,
		initial: initialValue,
	}
}

func (cg *CongruentialGenerator) Name() string {
	return string(cg.name)
}

func (cg *CongruentialGenerator) String() string {
	d := make(map[string]interface{}, 5)

	d["distributionName"] = cg.name
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
	name GeneratorName
	g    IntGenerator
	m    int
}

func (ug *UniformGenerator) Name() string {
	return string(ug.name)
}

func (ug *UniformGenerator) String() string {
	d := make(map[string]interface{}, 2)

	d["distributionName"] = ug.name
	d["modulus"] = ug.m

	if b, err := json.Marshal(d); err != nil {
		return fmt.Sprintf("%v", d)
	} else {
		return string(b)
	}
}

func NewUniformGeneratorDefault() *UniformGenerator {
	return &UniformGenerator{name: Uniform, g: rand.New(rand.NewSource(rand.Int63())), m: math.MaxInt32}
}

func NewUniformGenerator(generator IntGenerator, modulus int) *UniformGenerator {
	return &UniformGenerator{
		name: Uniform,
		g:    generator,
		m:    modulus,
	}
}

func (ug *UniformGenerator) Float64() float64 {
	if ug.m == 0 {
		panic("modulus is 0")
	}

	return float64(ug.g.Int()) / float64(ug.m)
}

type ExponentialGenerator struct {
	name GeneratorName
	g    Float64Generator
	l    float64
}

func (eg *ExponentialGenerator) Name() string {
	return string(eg.name)
}

func (eg *ExponentialGenerator) String() string {
	d := make(map[string]interface{}, 2)

	d["distributionName"] = eg.name
	d["rate"] = eg.l

	if b, err := json.Marshal(d); err != nil {
		return fmt.Sprintf("%v", d)
	} else {
		return string(b)
	}
}

func NewExponentialGeneratorDefault() *ExponentialGenerator {
	return &ExponentialGenerator{name: Exponential, g: rand.New(rand.NewSource(rand.Int63())), l: 1}
}

func NewExponentialGenerator(generator Float64Generator, rate float64) *ExponentialGenerator {
	return &ExponentialGenerator{name: Exponential, g: generator, l: rate}
}

func (eg *ExponentialGenerator) ExpFloat64() float64 {
	return -eg.l * math.Log(eg.g.Float64())
}

type NormalGenerator struct {
	name   GeneratorName
	g      Float64Generator
	g2     Float64Generator
	stdDev float64
	mean   float64
}

func (ng *NormalGenerator) Name() string {
	return string(ng.name)
}

func (ng *NormalGenerator) String() string {
	d := make(map[string]interface{}, 3)

	d["distributionName"] = ng.name
	d["standardDeviation"] = ng.stdDev
	d["mean"] = ng.mean

	if b, err := json.Marshal(d); err != nil {
		return fmt.Sprintf("%v", d)
	} else {
		return string(b)
	}
}

func NewNormalGeneratorDefault() *NormalGenerator {
	return &NormalGenerator{name: Normal, g: rand.New(rand.NewSource(rand.Int63())), g2: rand.New(rand.NewSource(rand.Int63())), stdDev: 1, mean: 0}
}

func NewNormalGenerator(generator Float64Generator, secondGenerator Float64Generator, standardDeviation float64, mean float64) *NormalGenerator {
	return &NormalGenerator{name: Normal, g: generator, g2: secondGenerator, stdDev: standardDeviation, mean: mean}
}

func (ng *NormalGenerator) NormFloat64() float64 {
	v1 := 2*ng.g.Float64() - 1
	v2 := 2*ng.g2.Float64() - 1
	S := math.Pow(v1, 2) + math.Pow(v2, 2)

	if S >= 1 {
		return ng.NormFloat64()
	}

	intermediate := math.Sqrt(-2/S*math.Log(S)) * ng.stdDev

	return intermediate*v1 + ng.mean
}

type FloatPair struct {
	x float64
	y float64
}

type FloatPairs []FloatPair

func (p FloatPairs) Len() int {
	return len(p)
}

func (p FloatPairs) XY(idx int) (x, y float64) {
	return p[idx].x, p[idx].y
}

type TwoDimensionalGenerator struct {
	name GeneratorName

	g  Float64Generator
	g2 Float64Generator

	// sigma_x
	stdDevX float64
	// sigma_y
	stdDevY float64
	// m_x
	meanX float64
	// m_y
	meanY float64
	// r
	correlationCoefficient float64
}

func (tdg *TwoDimensionalGenerator) Name() string {
	return string(tdg.name)
}

func (tdg *TwoDimensionalGenerator) String() string {
	d := make(map[string]interface{}, 6)

	d["distributionName"] = tdg.name
	d["standardDeviationX"] = tdg.stdDevX
	d["standardDeviationY"] = tdg.stdDevY
	d["meanX"] = tdg.meanX
	d["meanY"] = tdg.meanY
	d["correlationCoefficient"] = tdg.correlationCoefficient

	if b, err := json.Marshal(d); err != nil {
		return fmt.Sprintf("%v", d)
	} else {
		return string(b)
	}
}

func NewTwoDimensionalGeneratorDefault() *TwoDimensionalGenerator {
	return &TwoDimensionalGenerator{
		name:                   TwoDimensional,
		g:                      rand.New(rand.NewSource(rand.Int63())),
		g2:                     rand.New(rand.NewSource(rand.Int63())),
		stdDevX:                1,
		stdDevY:                1,
		meanX:                  0,
		meanY:                  0,
		correlationCoefficient: 0.5,
	}
}

func NewTwoDimensionalGenerator(
	generator Float64Generator,
	secondGenerator Float64Generator,
	standardDeviationX float64,
	standardDeviationY float64,
	meanX float64,
	meanY float64,
	correlationCoefficient float64,
) *TwoDimensionalGenerator {
	return &TwoDimensionalGenerator{
		name:                   TwoDimensional,
		g:                      generator,
		g2:                     secondGenerator,
		stdDevX:                standardDeviationX,
		stdDevY:                standardDeviationY,
		meanX:                  meanX,
		meanY:                  meanY,
		correlationCoefficient: correlationCoefficient,
	}
}

func (tdg *TwoDimensionalGenerator) TwoDimensionalFloat64s() FloatPair {
	SxComponents := make([]float64, 0, 6)
	SyComponents := make([]float64, 0, 6)

	for i := 0; i < 6; i += 1 {
		SxComponents = append(SxComponents, tdg.g.Float64())
		SyComponents = append(SyComponents, tdg.g2.Float64())
	}

	Sx := floats.Sum(SxComponents)
	Sy := floats.Sum(SyComponents)

	x := math.Sqrt2*tdg.stdDevX*(Sx-3) + tdg.meanX

	meanYX := tdg.meanY + tdg.correlationCoefficient*(x-tdg.meanX)*tdg.stdDevY/tdg.stdDevX
	stdDevYX := tdg.stdDevY * math.Sqrt(1-math.Pow(tdg.correlationCoefficient, 2))

	y := math.Sqrt2*stdDevYX*(Sy-3) + meanYX

	return FloatPair{x: x, y: y}
}
