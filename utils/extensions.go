package utils

func Contains[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

func All[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}
	return true
}

func SplitToBatches[T any](input []T, batchSize int) [][]T {
	var batches [][]T
	total := len(input)

	if input == nil || total == 0 {
		return [][]T{}
	}

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batches = append(batches, input[i:end])
	}

	return batches
}
