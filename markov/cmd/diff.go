package main

type info struct {
	t    float64
	e    float64
	diff float64
}

func diff(tp, ep []float64) []info {
	d := make([]info, 0, len(tp))
	for i := 0; i < len(tp); i++ {
		d = append(d, info{
			t:    tp[i],
			e:    ep[i],
			diff: ep[i] - tp[i],
		})
	}

	return d
}
