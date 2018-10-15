package mergesort

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const size = 100000

func TestMergesort(t *testing.T) {
	s := []int{5, 8, 9, 5, 0, 10, 1, 6}
	mergesort(s)
	assert.Equal(t, []int{0, 1, 5, 5, 6, 8, 9, 10}, s)
}

func TestParallelMergesort(t *testing.T) {
	s := []int{5, 8, 9, 5, 0, 10, 1, 6}
	parallelMergesort(s)
	assert.Equal(t, []int{0, 1, 5, 5, 6, 8, 9, 10}, s)
}

func TestParallelMergesortWithChannel(t *testing.T) {
	s := []int{5, 8, 9, 5, 0, 10, 1, 6}
	ch := initParallelMergesortWithChannel()
	parallelMergesortWithChannel(ch, s)
	assert.Equal(t, []int{0, 1, 5, 5, 6, 8, 9, 10}, s)
}

func BenchmarkMergesort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := random(size)
		b.StartTimer()
		mergesort(s)
		b.StopTimer()
	}
}

func BenchmarkParallelMergesort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := random(size)
		b.StartTimer()
		parallelMergesort(s)
		b.StopTimer()
	}
}

func BenchmarkParallelMergesortWithChannel(b *testing.B) {
	ch := initParallelMergesortWithChannel()
	defer close(ch)

	for i := 0; i < b.N; i++ {
		s := random(size)
		parallelMergesortWithChannel(ch, s)
	}
}
