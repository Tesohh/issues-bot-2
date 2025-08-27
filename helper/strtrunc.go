package helper

func StrTrunc(str string, maxLen int) string {
	if maxLen == 0 {
		return str
	}

	cut := str
	if len(str) > maxLen {
		cut = cut[:maxLen-1]
		cut += "…"
	}

	return cut
}
