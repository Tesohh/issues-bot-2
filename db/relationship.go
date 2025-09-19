package db

import "time"

type RelationshipKind string

const (
	RelationshipKindDependency RelationshipKind = "dependency"
)

type Relationship struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	FromIssueID uint
	FromIssue   Issue

	ToIssueID uint
	ToIssue   Issue

	Kind RelationshipKind `gorm:"check:kind in ('dependency')"`
}

// helper type used by some functions to group relationships by their direction
type RelationshipsByDirection struct {
	Inbound  []Relationship
	Outbound []Relationship
}
