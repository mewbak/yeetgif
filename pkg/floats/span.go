package floats

import "math"

func MakeSpan(min, max float64, nPoints int) []float64 {
	X := make([]float64, nPoints)
	min, max = math.Min(max, min), math.Max(max, min)
	d := max - min
	for i := range X {
		X[i] = min + d*(float64(i)/float64(nPoints-1))
	}
	return X
}

func Span(min, max float64, X []float64) {
	min, max = math.Min(max, min), math.Max(max, min)
	d := max - min
	for i := range X {
		X[i] = min + d*(float64(i)/float64(len(X)-1))
	}
}
