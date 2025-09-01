package helper

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
