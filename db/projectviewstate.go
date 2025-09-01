package db

import (
	"time"
)

// used for keeping information on the autolist display
type ProjectViewState struct {
	MessageID string `gorm:"primarykey"`
	ChannelID string
	CreatedAt time.Time
	UpdatedAt time.Time

	ProjectID uint
	Project   Project

	CurrentPage int

	ListNameFmt string // eg. "# Autolist for %s `[%s]`"

	Filter IssueFilter `gorm:"type:json"`
	Sorter IssueSorter `gorm:"type:json"`
}
