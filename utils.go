package golok

func insert[T any](s []T, index int, value T) []T {
	if index < 0 || index > len(s) {
		panic("index out of range")
	}
	s = append(s[:index], append([]T{value}, s[index:]...)...)
	return s
}

func remove[T any](s []T, index int) []T {
	if index < 0 || index >= len(s) {
		panic("index out of range")
	}
	return append(s[:index], s[index+1:]...)
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
