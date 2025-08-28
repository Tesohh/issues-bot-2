package db

import (
	"time"
)

// used for keeping information on the autolist display
type ProjectViewState struct {
	MessageID string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ProjectID uint
	Project   Project

	CurrentPage int

	Filter IssueFilter `gorm:"type:json"`
	Sorter IssueSorter `gorm:"type:json"`
}
