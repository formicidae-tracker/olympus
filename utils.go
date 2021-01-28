package main

// searches from index n-1 to 0 the first index which is true or -1 if none are true

func BackLinearSearch(n int, test func(int) bool) int {
	for n > 0 {
		n--
		if test(n) == true {
			return n
		}
	}
	return n - 1
}

func LinearSearch(n int, test func(int) bool) int {

	for i := 0; i < n; i++ {
		if test(i) == true {
			return i
		}
	}
	return n
}
