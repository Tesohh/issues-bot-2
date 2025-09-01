package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type IssueSortBy string

const (
	IssueSortByCode IssueSortBy = "code"
	IssueSortByDate IssueSortBy = "date"
)

type SortOrder int

const (
	SortOrderAscending  SortOrder = 0
	SortOrderDescending SortOrder = 1
)

type IssueSorter struct {
	SortBy    IssueSortBy
	SortOrder SortOrder
}

func (sorter IssueSorter) Apply(issues []Issue) []Issue {
	// TODO: implement
	return issues
}

func (sorter IssueSorter) String() string {
	// TODO: implement
	return ""
}

// functions to save this in the db as json
func (s IssueSorter) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *IssueSorter) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan IssueFilter")
	}
	return json.Unmarshal(bytes, s)
}
