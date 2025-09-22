package helper

func PagesWithOffset[T any](slice []T, pageSize int, offset int) int {
	if pageSize <= 0 {
		return 0
	}
	n := len(slice)

	if n <= offset {
		return 1
	}

	remaining := n - offset
	fullPages := (remaining + pageSize - 1) / pageSize

	return 1 + fullPages
}

func PaginateWithOffset[T any](slice []T, pageSize int, page int, offset int) []T {
	if page < 0 {
		page = 0
	}

	var start, end int
	if page == 0 {
		start = 0
		end = offset + pageSize
	} else {
		start = offset + (page-1)*pageSize
		end = start + pageSize
	}

	if start > len(slice) {
		return []T{}
	}
	if end > len(slice) {
		end = len(slice)
	}

	return slice[start:end]
}
