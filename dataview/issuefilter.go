package dataview

import "issues/v2/db"

type IssueFilter struct {
}

func (filter IssueFilter) Apply(issues []db.Issue) []db.Issue {
	return issues
}

func (filter IssueFilter) String() string {
	return ""
}
