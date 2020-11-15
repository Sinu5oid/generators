package stat

import "math"

type interval struct {
	values     []float64
	leftBound  float64
	rightBound float64
}

func filterValues(source []float64, leftBound float64, rightBound float64, includeRightBound bool) []float64 {
	result := make([]float64, 0)
	for _, v := range source {
		if v < leftBound || v > rightBound {
			continue
		}

		if !includeRightBound && v == rightBound {
			continue
		}

		result = append(result, v)
	}

	return result
}

func split(source []float64, bounds []float64) []interval {
	if len(bounds) < 2 {
		panic("not enough bounds")
	}

	result := make([]interval, 0, len(bounds))

	for i := 1; i < len(bounds); i += 1 {
		result = append(result, interval{
			values:     filterValues(source, bounds[i-1], bounds[i], i == len(bounds)-1),
			leftBound:  bounds[i-1],
			rightBound: bounds[i],
		})
	}

	return result
}

func merge(source []interval) []interval {
	result := make([]interval, 0, cap(source))

	buffer := interval{}
	for i, group := range source {
		if len(group.values)+len(buffer.values) < 5 {
			// less than 5 items given that there are items in buffer
			// (last remembered)
			if i < len(source)-1 {
				// not the last group
				buffer.leftBound = group.leftBound
				buffer.values = append(buffer.values, group.values...)
			} else {
				// last group (end of slice)

				// take previous group
				buffer.values = append(buffer.values, result[len(result)-1].values...)
				// append with current group
				buffer.leftBound = result[len(result)-1].leftBound
				buffer.values = append(buffer.values, group.values...)
				buffer.rightBound = group.rightBound
				// replace previous group
				result = append(result[:len(result)-1], buffer)

				// clear buffer
				buffer = interval{}
			}
		} else {
			if len(buffer.values) == 0 {
				result = append(result, interval{
					values:     group.values,
					leftBound:  group.leftBound,
					rightBound: group.rightBound,
				})
			} else {
				result = append(result, interval{
					values:     append(buffer.values, group.values...),
					leftBound:  buffer.leftBound,
					rightBound: group.rightBound,
				})
			}

			buffer = interval{}
		}
	}

	return result
}

func normalCDF(x float64, sigma float64, alpha float64) float64 {
	return (1 + math.Erf((x-alpha)/(sigma*math.Sqrt2))) / 2
}

func exponentialCDF(x float64, lambda float64) float64 {
	if x < 0 {
		return 0
	}

	return 1 - math.Pow(math.E, -lambda*x)
}

func uniformCDF(x float64, a float64, b float64) float64 {
	if x <= a {
		return 0
	}

	if x > b {
		return 1
	}

	return (x - a) / (b - a)
}
