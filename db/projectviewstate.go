package db

import (
	"time"

	"gorm.io/gorm"
)

// used for keeping information on the autolist display
type ProjectViewState struct {
	MessageID string `gorm:"primarykey"`
	ChannelID string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ProjectID uint
	Project   Project

	Permanent bool

	CurrentPage int

	ListNameFmt string // eg. "# Autolist for %s `[%s]`"

	Filter IssueFilter `gorm:"type:json"`
	Sorter IssueSorter `gorm:"type:json"`
}
