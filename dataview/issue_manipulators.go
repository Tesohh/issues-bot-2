package dataview

import "issues/v2/db"

type IssueFilter struct {
	Statuses        []db.IssueStatus
	Tags            []string
	PriorityRoleIDs []string
	CategoryRoleIDs []string
	AssigneeIDs     []string
}

func (filter IssueFilter) Apply(issues []db.Issue) []db.Issue {
	// TODO: implement
	return issues
}

func (filter IssueFilter) String() string {
	// TODO: implement
	return ""
}

type IssueSortBy uint8

const (
	IssueSortByCode IssueSortBy = 0
	IssueSortByDate IssueSortBy = 1
)

type SortOrder uint8

const (
	SortOrderAscending  SortOrder = 0
	SortOrderDescending SortOrder = 1
)

type IssueSorter struct {
	SortBy    IssueSortBy
	SortOrder SortOrder
}

func (sorter IssueSorter) Apply(issues []db.Issue) []db.Issue {
	// TODO: implement
	return issues
}

func (sorter IssueSorter) String() string {
	// TODO: implement
	return ""
}
