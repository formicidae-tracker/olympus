package main

// searches from index n-1 to 0 the first index which is true or -1 if none are true

func Insert[T any](slice []T, value T, index int) []T {
	res := append(slice, value)
	copy(res[index+1:], res[index:])
	res[index] = value
	return res
}

func BackInsertionSort[T any](slice []T, value T, less func(a, b T) bool) []T {
	i := BackLinearSearch(slice, value, less)
	return Insert(slice, value, i)
}

func InsertionSort[T any](slice []T, value T, less func(a, b T) bool) []T {
	i := LinearSearch(slice, value, less)
	return Insert(slice, value, i)
}

func BackLinearSearch[T any](slice []T, value T, less func(a, b T) bool) int {
	n := len(slice)
	for n > 0 {
		n--
		if less(value, slice[n]) == true {
			continue
		}
		return n + 1
	}
	return 0
}

func LinearSearch[T any](slice []T, value T, less func(a, b T) bool) int {

	for i := 0; i < len(slice); i++ {
		if less(value, slice[i]) == true {
			return i
		}
	}
	return len(slice)
}

func ZoneIdentifier(hostname, zoneName string) string {
	return hostname + "." + zoneName
}
