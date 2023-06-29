package weight_random

type WItem interface {
	Weight() uint64
}

type WeightSlice[T WItem] []T

// Rand Binomial Distribution
func (x WeightSlice[T]) Rand(r uint64) T {
	var zero T
	if len(x) == 0 {
		return zero
	}
	var total uint64 = 0
	for i := 0; i < len(x); i++ {
		total += x[i].Weight()
	}
	if total <= 0 {
		return zero
	}
	cursor, val := uint64(0), r%total
	for i := 0; i < len(x); i++ {
		cursor += x[i].Weight()
		if val < cursor {
			return x[i]
		}
	}
	return zero
}
