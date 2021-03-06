package generators

import (
	"math"
	"math/rand"
	"testing"
)

func BenchmarkCongruentialGenerator(b *testing.B) {
	g := NewCongruentialGenerator(int(math.Pow(2, 32)), 1103515245, 12345, 0)

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		_ = g.Int()
	}
}

func BenchmarkStdCongruentialGenerator(b *testing.B) {
	for i := 0; i < b.N; i += 1 {
		_ = rand.Int()
	}
}

func BenchmarkUniformGenerator(b *testing.B) {
	s := rand.NewSource(rand.Int63())

	g := NewUniformGenerator(rand.New(s), 1)

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		_ = g.Float64()
	}
}

func BenchmarkUniformGeneratorDefault(b *testing.B) {
	g := NewUniformGeneratorDefault()

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		_ = g.Float64()
	}
}

func BenchmarkStdUniformGenerator(b *testing.B) {
	for i := 0; i < b.N; i += 1 {
		_ = rand.Float64()
	}
}

func BenchmarkChainedUniformGenerator(b *testing.B) {
	cg := NewCongruentialGenerator(int(math.Pow(2, 32)), 1103515245, 12345, 0)
	ug := NewUniformGenerator(cg, int(math.Pow(2, 32)))

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		_ = ug.Float64()
	}
}

func BenchmarkExponentialGenerator(b *testing.B) {
	s := rand.NewSource(rand.Int63())

	g := NewExponentialGenerator(rand.New(s), 1)

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		_ = g.ExpFloat64()
	}
}

func BenchmarkExponentialGeneratorDefault(b *testing.B) {
	g := NewExponentialGeneratorDefault()

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		_ = g.ExpFloat64()
	}
}

func BenchmarkStdExponentialGenerator(b *testing.B) {
	for i := 0; i < b.N; i += 1 {
		_ = rand.ExpFloat64()
	}
}

func BenchmarkChainedExponentialGenerator(b *testing.B) {
	cg := NewCongruentialGenerator(int(math.Pow(2, 32)), 1103515245, 12345, 0)
	ug := NewUniformGenerator(cg, int(math.Pow(2, 32)))
	eg := NewExponentialGenerator(ug, 345)

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		_ = eg.ExpFloat64()
	}
}

func BenchmarkNormalGenerator(b *testing.B) {
	s := rand.NewSource(rand.Int63())
	s2 := rand.NewSource(rand.Int63())

	g := NewNormalGenerator(rand.New(s), rand.New(s2), 1, 0)

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		_ = g.NormFloat64()
	}
}

func BenchmarkNormalGeneratorDefault(b *testing.B) {
	g := NewNormalGeneratorDefault()

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		_ = g.NormFloat64()
	}
}

func BenchmarkStdNormalGenerator(b *testing.B) {
	for i := 0; i < b.N; i += 1 {
		_ = rand.NormFloat64()
	}
}

func BenchmarkChainedNormalGenerator(b *testing.B) {
	cg := NewCongruentialGenerator(int(math.Pow(2, 32)), 1103515245, 12345, 0)
	ug := NewUniformGenerator(cg, int(math.Pow(2, 32)))

	cg2 := NewCongruentialGenerator(int(math.Pow(2, 32)), 134775813, 1, 3)
	ug2 := NewUniformGenerator(cg2, int(math.Pow(2, 32)))

	ng := NewNormalGenerator(ug, ug2, 1, 2)

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		_ = ng.NormFloat64()
	}
}

func BenchmarkTwoDimensionalGenerator(b *testing.B) {
	s := rand.NewSource(rand.Int63())
	s2 := rand.NewSource(rand.Int63())

	tdg := NewTwoDimensionalGenerator(rand.New(s), rand.New(s2), 1, 1, 0, 0, 0.5)

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		_ = tdg.TwoDimensionalFloat64s()
	}
}

func BenchmarkChainedTwoDimensionalGenerator(b *testing.B) {
	cg := NewCongruentialGenerator(int(math.Pow(2, 32)), 1103515245, 12345, 0)
	ug := NewUniformGenerator(cg, int(math.Pow(2, 32)))

	cg2 := NewCongruentialGenerator(int(math.Pow(2, 32)), 134775813, 1, 3)
	ug2 := NewUniformGenerator(cg2, int(math.Pow(2, 32)))

	tdg := NewTwoDimensionalGenerator(ug, ug2, 1, 1, 0, 0, 0.5)

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		_ = tdg.TwoDimensionalFloat64s()
	}
}

func BenchmarkTwoDimensionalGeneratorDefault(b *testing.B) {
	tdg := NewTwoDimensionalGeneratorDefault()

	b.ResetTimer()

	for i := 0; i < b.N; i += 1 {
		_ = tdg.TwoDimensionalFloat64s()
	}
}
