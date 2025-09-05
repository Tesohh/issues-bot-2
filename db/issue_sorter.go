package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"slices"
)

type IssueSortBy string

const (
	IssueSortByCode IssueSortBy = "code"
	IssueSortByDate IssueSortBy = "date"
)

type SortOrder string

const (
	SortOrderAscending  SortOrder = "asc"
	SortOrderDescending SortOrder = "desc"
)

type IssueSorter struct {
	SortBy    IssueSortBy
	SortOrder SortOrder
}

func DefaultSorter() IssueSorter {
	return IssueSorter{
		SortBy:    IssueSortByCode,
		SortOrder: SortOrderAscending,
	}
}

func (sorter IssueSorter) Apply(issues []Issue) []Issue {
	return slices.SortedFunc(slices.Values(issues), func(a, b Issue) int {
		if sorter.SortOrder == SortOrderDescending {
			a, b = b, a
		}

		switch sorter.SortBy {
		case IssueSortByDate:
			return a.UpdatedAt.Compare(b.UpdatedAt)
		case IssueSortByCode:
			if a.Code == nil || b.Code == nil {
				return 0
			}
			return int(*a.Code) - int(*b.Code)
		}

		return 0
	})
}

func (sorter IssueSorter) String() string {
	str := "sort by"
	str += " " + string(sorter.SortBy)
	switch sorter.SortOrder {
	case SortOrderAscending:
		str += " asc"
	case SortOrderDescending:
		str += " desc"
	}

	return str
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
