package db

type Project struct {
	ID      uint `gorm:"primarykey"`
	Name    string
	Prefix  string
	RepoURL string

	DiscordCategoryChannelID string
	IssuesInputChannelID     string
	AutoListMessageID        string

	GuildID string

	Issues []Issue
}
