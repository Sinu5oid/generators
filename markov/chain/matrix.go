package chain

import "github.com/pkg/errors"

var (
	ErrInvalidMatrix    = errors.New("invalid transition matrix")
	ErrInvalidArguments = errors.New("invalid arguments")
)

const tol = 1e-9

func validateMatrix(tm [][]float64) error {
	if tm == nil {
		return errors.Wrap(ErrInvalidMatrix, "nil")
	}

	if len(tm) == 0 {
		return errors.Wrap(ErrInvalidMatrix, "empty")
	}

	ref := len(tm[0])
	for i := 0; i < len(tm); i++ {
		if ref != len(tm[i]) {
			return errors.Wrap(ErrInvalidMatrix, "inconsistent length")
		}

		sumP := float64(0)
		for j := 0; j < ref; j++ {
			if tm[i][j] < 0 {
				return errors.Wrap(ErrInvalidMatrix, "contains negative items")
			}

			sumP += tm[i][j]
		}

		if sumP > 1+tol {
			return errors.Wrap(ErrInvalidMatrix, "sum of row is greater than 1")
		}
	}

	return nil
}

func mustMultiplyMatrices(a, b [][]float64) [][]float64 {
	if len(a[0]) != len(b) {
		panic(errors.Wrap(ErrInvalidArguments, "lengths are not comparable"))
	}

	res := make([][]float64, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = make([]float64, len(b[0]))
		for j := 0; j < len(b[0]); j++ {
			for k := 0; k < len(b); k++ {
				res[i][j] += a[i][k] * b[k][j]
			}
		}
	}

	return res
}
