package helper

func Pages[T any](slice []T, pageSize int) int {
	return (len(slice) / pageSize) + 1
}

func Paginate[T any](slice []T, pageSize int, page int) []T {
	if page < 0 {
		page = 0
	}

	start := page * pageSize
	if len(slice) < start {
		return []T{}
	}

	end := start + pageSize
	if len(slice) < end {
		end = len(slice)
	}

	return slice[start:end]
}
