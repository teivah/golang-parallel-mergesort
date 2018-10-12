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

/* Sequential */

func mergesort(s []int) {
	if len(s) > 1 {
		middle := len(s) / 2
		mergesort(s[:middle])
		mergesort(s[middle:])
		merge(s, middle)
	}
}

/* Parallel */

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

/* Parallel with channels */

type Data struct {
	s  []int
	wg *sync.WaitGroup
}

func initParallelMergesortWithChannel() chan Data {
	cores := runtime.NumCPU()

	ch := make(chan Data, 1024)

	for i := 0; i < cores; i++ {
		go parallelMergesortWithChannelHandler(ch)
	}

	return ch
}

func parallelMergesortWithChannel(ch chan Data, s []int) {
	var wg sync.WaitGroup
	wg.Add(1)

	ch <- Data{
		s:  s,
		wg: &wg,
	}

	wg.Wait()
}

func parallelMergesortWithChannelHandler(ch chan Data) {
	for data := range ch {
		s := data.s
		len := len(s)

		if len > 1 {
			if len <= max { // Sequential
				mergesort(s)
			} else { // Parallel
				middle := len / 2

				var wg sync.WaitGroup
				wg.Add(2)

				ch <- Data{
					s:  s[:middle],
					wg: &wg,
				}

				ch <- Data{
					s:  s[middle:],
					wg: &wg,
				}

				wg.Wait()
				merge(s, middle)
			}
		}

		data.wg.Done()
	}
}
