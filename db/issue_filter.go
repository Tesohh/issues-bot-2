package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type IssueFilter struct {
	Statuses        []IssueStatus
	Title           string
	Tags            []string
	PriorityRoleIDs []string
	CategoryRoleIDs []string
	AssigneeIDs     []string
}

func (s IssueFilter) Apply(issues []Issue) []Issue {
	// TODO: implement
	return issues
}

func (s IssueFilter) String() string {
	// TODO: implement
	return ""
}

// functions to save this in the db as json
func (s IssueFilter) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *IssueFilter) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan IssueFilter")
	}
	return json.Unmarshal(bytes, s)
}
