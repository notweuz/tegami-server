package util

func MapToList[M, T any](models []M, fn func(M) T) []T {
	items := make([]T, len(models))
	for i, m := range models {
		items[i] = fn(m)
	}
	return items
}