package util

func Last[T any](list []T) *T {
	back := Back(list)
	if back < 0 {
		return nil
	}
	return &list[back]
}

func Back[T any](list []T) int {
	return len(list) - 1
}
