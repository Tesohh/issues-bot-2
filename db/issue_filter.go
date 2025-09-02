package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

var issueStatusNames = [4]string{"todo", "working", "done", "killed"} // TODO: don't duplicate this (for now we have import cycle)

type IssueFilter struct {
	Statuses        []IssueStatus
	Tags            []string
	PriorityRoleIDs []string
	CategoryRoleIDs []string
	AssigneeIDs     []string
	Title           string
}

func DefaultFilter() IssueFilter {
	return IssueFilter{
		Statuses:        []IssueStatus{},
		Tags:            []string{},
		PriorityRoleIDs: []string{},
		CategoryRoleIDs: []string{},
		AssigneeIDs:     []string{},
		Title:           "",
	}
}

func (f IssueFilter) Apply(issues []Issue) []Issue {
	return issues
}

func (f IssueFilter) String() string {
	keywords := []string{}
	for _, status := range f.Statuses {
		keywords = append(keywords, issueStatusNames[status])
	}
	for _, tag := range f.Tags {
		keywords = append(keywords, fmt.Sprintf("`+%s`", tag))
	}
	for _, roleID := range append(f.PriorityRoleIDs, f.CategoryRoleIDs...) {
		keywords = append(keywords, fmt.Sprintf("<@&%s>", roleID))
	}
	for _, assigneeID := range f.AssigneeIDs {
		keywords = append(keywords, fmt.Sprintf("<@%s>", assigneeID))
	}
	if len(f.Title) > 0 {
		keywords = append(keywords, fmt.Sprintf("\"%s\"", f.Title))
	}

	if len(keywords) == 0 {
		keywords = append(keywords, "all")
	}

	return strings.Join(keywords, ", ")
}

// functions to save this in the db as json
func (f IssueFilter) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *IssueFilter) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan IssueFilter")
	}
	return json.Unmarshal(bytes, f)
}
