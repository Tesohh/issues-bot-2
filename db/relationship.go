package db

type RelationshipKind string

const (
	RelationshipKindDependency RelationshipKind = ""
)

type Relationship struct {
	ID uint `gorm:"primarykey"`

	FromIssueID uint
	FromIssue   Issue

	ToIssueID uint
	ToIssue   Issue

	Kind RelationshipKind
}
