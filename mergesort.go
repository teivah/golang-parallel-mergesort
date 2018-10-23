package mergesort

import (
	"sync"
)

const max = 1 << 11

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

/* Parallel 1 */

func parallelMergesort1(s []int) {
	parallelMergesortHandler1(s, nil)
}

func parallelMergesortHandler1(s []int, parent *sync.WaitGroup) {
	len := len(s)

	if len > 1 {
		if len <= max { // Sequential
			mergesort(s)
		} else { // Parallel
			middle := len / 2

			var wg sync.WaitGroup
			wg.Add(2)

			go parallelMergesortHandler1(s[:middle], &wg)
			go parallelMergesortHandler1(s[middle:], &wg)

			wg.Wait()
			merge(s, middle)
		}
	}

	if parent != nil {
		parent.Done()
	}
}

/* Parallel 2 */

func parallelMergesort2(s []int) {
	parallelMergesortHandler2(s, nil)
}

func parallelMergesortHandler2(s []int, parent *sync.WaitGroup) {
	len := len(s)

	if len > 1 {
		if len <= max { // Sequential
			mergesort(s)
		} else { // Parallel
			middle := len / 2

			var wg sync.WaitGroup
			wg.Add(1)

			go parallelMergesortHandler2(s[:middle], &wg)
			parallelMergesortHandler2(s[middle:], nil)

			wg.Wait()
			merge(s, middle)
		}
	}

	if parent != nil {
		parent.Done()
	}
}

/* Parallel 3 */

func parallelMergesort3(s []int) {
	parallelMergesortHandler3(s, nil)
}

func parallelMergesortHandler3(s []int, parent *sync.WaitGroup) {
	len := len(s)

	if len > 1 {
		middle := len / 2

		var wg sync.WaitGroup
		wg.Add(2)

		go parallelMergesortHandler3(s[:middle], &wg)
		go parallelMergesortHandler3(s[middle:], &wg)

		wg.Wait()
		merge(s, middle)
	}

	if parent != nil {
		parent.Done()
	}
}
