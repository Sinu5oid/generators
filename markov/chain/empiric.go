package chain

func EProb(impls [][]int, tml, t, ic int) []float64 {
	result := make([]float64, tml, tml)

	for impl := 0; impl < ic; impl++ {
		if len(impls[impl]) <= t {
			result[impls[impl][len(impls[impl])-1]]++
			continue
		}

		result[impls[impl][t]]++
	}

	for i := 0; i < tml; i++ {
		result[i] /= float64(ic)
	}

	return result
}
