package weight_random

import (
	"math"
	"math/rand"
	"testing"
)

type item struct {
	weight uint64
}

func (i item) Weight() uint64 {
	return i.weight
}

func TestWeightSlice_Rand(t *testing.T) {
	sl := WeightSlice[*item]{}
	v := sl.Rand(rand.Uint64())
	if v != nil {
		t.Error()
	}
	sl = append(sl, &item{weight: 0})
	v = sl.Rand(rand.Uint64())
	if v != nil {
		t.Error()
	}
	sl = append(sl, &item{weight: 1})
	v = sl.Rand(rand.Uint64())
	if v == nil || v.weight != 1 {
		t.Error()
	}
	sl = append(sl, &item{weight: math.MaxUint64 - 1})
	v = sl.Rand(rand.Uint64())
	if v == nil || v.weight != math.MaxUint64-1 {
		t.Error()
	}
}
