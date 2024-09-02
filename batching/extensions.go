package batching

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
