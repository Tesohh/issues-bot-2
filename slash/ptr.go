package slash

// Helps making values for discordgo
func Ptr[T any](v T) *T {
	return &v
}
