package db

import "time"

type RelationshipKind string

const (
	RelationshipKindDependency RelationshipKind = ""
)

type Relationship struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	FromIssueID uint
	FromIssue   Issue

	ToIssueID uint
	ToIssue   Issue

	Kind RelationshipKind `gorm:"check:kind in ('')"`
}
