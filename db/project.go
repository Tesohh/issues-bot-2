package db

import "time"

type Project struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name    string
	Prefix  string
	RepoURL string

	DiscordCategoryChannelID string
	GeneralChannelID         string
	IssuesInputChannelID     string
	AutoListMessageID        string

	HasBeenSolicitedByListWarning bool

	GuildID string

	ProjectViewState []ProjectViewState
	Issues           []Issue
}
