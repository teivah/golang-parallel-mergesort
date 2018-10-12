package mergesort

import (
	"runtime"
	"sync"
)

const max = 1 << 13

func merge(s []int, middle int) {
	helper := make([]int, len(s))
	copy(helper, s)

	helperLeft := 0
	helperRight := middle
	current := 0
	high := len(s) - 1

	for helperLeft <= middle-1 && helperRight <= high {
		if helper[helperLeft] <= helper[helperRight] {
			s[current] = helper[helperLeft]
			helperLeft++
		} else {
			s[current] = helper[helperRight]
			helperRight++
		}
		current++
	}

	for helperLeft <= middle-1 {
		s[current] = helper[helperLeft]
		current++
		helperLeft++
	}
}

func mergesort(s []int) {
	if len(s) > 1 {
		middle := len(s) / 2
		mergesort(s[:middle])
		mergesort(s[middle:])
		merge(s, middle)
	}
}

func parallelMergesort(s []int) {
	parallelMergesortHandler(s, nil)
}

func parallelMergesortHandler(s []int, parent *sync.WaitGroup) {
	len := len(s)

	if len > 1 {
		if len <= max { // Sequential
			mergesort(s)
		} else { // Parallel
			middle := len / 2

			var wg sync.WaitGroup
			wg.Add(2)

			go parallelMergesortHandler(s[:middle], &wg)
			go parallelMergesortHandler(s[middle:], &wg)

			wg.Wait()
			merge(s, middle)
		}
	}

	if parent != nil {
		parent.Done()
	}
}

func parallelMergesortWithChannel(s []int) {
	cores := runtime.NumCPU()

	ch := make(chan []int)

	for i := 0; i < cores; i++ {
		go parallelMergesortWithChannelHandler(ch)
	}

	ch <- s
}

func parallelMergesortWithChannelHandler(ch chan []int) {
	for s := range ch {
		len := len(s)

		if len > 1 {
			middle := len/2 - 1

			ch <- s[:middle]
			ch <- s[middle+1:]
		}
	}
}
