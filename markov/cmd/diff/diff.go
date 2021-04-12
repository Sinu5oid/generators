package diff

type Info struct {
	T float64
	E float64
	D float64
}

func Get(tp, ep []float64) []Info {
	d := make([]Info, 0, len(tp))
	for i := 0; i < len(tp); i++ {
		d = append(d, Info{
			T: tp[i],
			E: ep[i],
			D: ep[i] - tp[i],
		})
	}

	return d
}
