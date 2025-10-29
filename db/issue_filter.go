package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

type IssueFilter struct {
	Statuses        []IssueStatus
	Tags            []string
	PriorityRoleIDs []string
	CategoryRoleIDs []string
	RecruiterIDs    []string
	AssigneeIDs     []string
	Nobody          bool
	Title           string
}

func DefaultFilter() IssueFilter {
	return IssueFilter{
		Statuses:        []IssueStatus{IssueStatusTodo, IssueStatusDoing},
		Tags:            []string{},
		PriorityRoleIDs: []string{},
		CategoryRoleIDs: []string{},
		RecruiterIDs:    []string{},
		AssigneeIDs:     []string{},
		Nobody:          false,
		Title:           "",
	}
}

func (f IssueFilter) isValid(issue Issue) bool {
	if len(f.Statuses) > 0 {
		if !slices.Contains(f.Statuses, issue.Status) {
			return false
		}
	}

	if len(f.Tags) > 0 {
		ok := false
		for _, tag := range issue.Tags {
			if slices.Contains(f.Tags, tag.Name) {
				ok = true
			}
		}
		if !ok {
			return false
		}
	}

	if len(f.PriorityRoleIDs) > 0 {
		if !slices.Contains(f.PriorityRoleIDs, issue.PriorityRoleID) {
			return false
		}
	}

	if len(f.CategoryRoleIDs) > 0 {
		if !slices.Contains(f.CategoryRoleIDs, issue.CategoryRoleID) {
			return false
		}
	}

	if len(f.RecruiterIDs) > 0 {
		if !slices.Contains(f.RecruiterIDs, issue.RecruiterUserID) {
			return false
		}
	}

	if f.Nobody {
		if len(issue.AssigneeUsers) != 0 {
			return false
		}
	} else if len(f.AssigneeIDs) > 0 {
		ok := false
		for _, assignee := range issue.AssigneeUsers {
			if slices.Contains(f.AssigneeIDs, assignee.ID) {
				ok = true
			}
		}
		if !ok {
			return false
		}
	}

	if len(f.Title) > 0 {
		if !strings.Contains(issue.Title, f.Title) {
			return false
		}
	}

	return true
}

func (f IssueFilter) Apply(issues []Issue) []Issue {
	filteredIssues := []Issue{}
	for _, issue := range issues {
		if f.isValid(issue) {
			filteredIssues = append(filteredIssues, issue)
		}
	}
	return filteredIssues
}

func (f IssueFilter) String() string {
	keywords := []string{}
	for _, status := range f.Statuses {
		keywords = append(keywords, IssueStatusNames[status])
	}
	for _, tag := range f.Tags {
		keywords = append(keywords, fmt.Sprintf("`+%s`", tag))
	}
	for _, roleID := range append(f.PriorityRoleIDs, f.CategoryRoleIDs...) {
		keywords = append(keywords, fmt.Sprintf("<@&%s>", roleID))
	}
	for _, recruiterID := range f.RecruiterIDs {
		keywords = append(keywords, fmt.Sprintf("<@%s>", recruiterID))
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
