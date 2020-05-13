package math2

import "math"

func Stddev(numbers []float64) *float64 {
	length := len(numbers)
	if length-1 <= 0 {
		return nil
	}

	var total, mean, sd float64
	for _, v := range numbers {
		total += v
	}
	mean = total / float64(length)

	for _, v := range numbers {
		sd += math.Pow(v-mean, 2)
	}

	v := math.Sqrt(sd / float64(length-1))
	return &v
}
