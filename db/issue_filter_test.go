package db_test

import (
	"issues/v2/db"
	"reflect"
	"testing"
)

var testSampleIssues = []db.Issue{
	{
		ID:              0,
		Title:           "Drink a gosser",
		Tags:            "gut,besser,gosser",
		Status:          db.IssueStatusTodo,
		RecruiterUserID: "tesohh",
		AssigneeUsers:   []db.User{{ID: "tesohh"}},
		CategoryRoleID:  "feat",
		PriorityRoleID:  "critical",
	},
	{
		ID:              1,
		Title:           "fix something",
		Tags:            "bug",
		Status:          db.IssueStatusWorking,
		RecruiterUserID: "lallos",
		AssigneeUsers:   []db.User{{ID: "tesohh"}},
		CategoryRoleID:  "fix",
		PriorityRoleID:  "normal",
	},
	{
		ID:              2,
		Title:           "SEND ADS",
		Tags:            "release,marketing",
		Status:          db.IssueStatusWorking,
		RecruiterUserID: "lallos",
		AssigneeUsers:   []db.User{{ID: "tesohh"}},
		CategoryRoleID:  "chore",
		PriorityRoleID:  "normal",
	},
	{
		ID:              3,
		Title:           "volantino ADS",
		Tags:            "release,marketing",
		Status:          db.IssueStatusWorking,
		RecruiterUserID: "tesohh",
		AssigneeUsers:   []db.User{{ID: "lallos"}},
		CategoryRoleID:  "chore",
		PriorityRoleID:  "important",
	},
	{
		ID:              4,
		Title:           "Finalize report",
		Tags:            "",
		Status:          db.IssueStatusDone,
		RecruiterUserID: "maria",
		AssigneeUsers:   []db.User{{ID: "maria"}, {ID: "tesohh"}},
		CategoryRoleID:  "docs",
		PriorityRoleID:  "low",
	},
	{
		ID:              5,
		Title:           "Deprecated feature cleanup",
		Tags:            "cleanup",
		Status:          db.IssueStatusKilled,
		RecruiterUserID: "admin",
		AssigneeUsers:   []db.User{{ID: "dev1"}, {ID: "dev2"}},
		CategoryRoleID:  "maintenance",
		PriorityRoleID:  "normal",
	},
	{
		ID:              6,
		Title:           "Extremely long issue title that should still be searchable by substring",
		Tags:            "performance, long title ",
		Status:          db.IssueStatusTodo,
		RecruiterUserID: "qa",
		AssigneeUsers:   []db.User{{ID: "dev1"}},
		CategoryRoleID:  "perf",
		PriorityRoleID:  "p0",
	},
	{
		ID:              7,
		Title:           "Hotfix urgent bug",
		Tags:            "bug, urgent",
		Status:          db.IssueStatusWorking,
		RecruiterUserID: "support",
		AssigneeUsers:   []db.User{},
		CategoryRoleID:  "urgent",
		PriorityRoleID:  "p0",
	},
	{
		ID:              8,
		Title:           "Marketing brainstorm",
		Tags:            "marketing,idea",
		Status:          db.IssueStatusTodo,
		RecruiterUserID: "marketing",
		AssigneeUsers:   []db.User{{ID: "creative"}},
		CategoryRoleID:  "brainstorm",
		PriorityRoleID:  "normal",
	},
	{
		ID:              9,
		Title:           "Design system revamp",
		Tags:            "design,refactor",
		Status:          db.IssueStatusWorking,
		RecruiterUserID: "uxlead",
		AssigneeUsers:   []db.User{{ID: "dev2"}},
		CategoryRoleID:  "ui",
		PriorityRoleID:  "critical",
	},
}

func TestIssueFilterApply(t *testing.T) {
	tests := []struct {
		name   string
		filter db.IssueFilter
		want   []int // expected issue IDs
	}{
		{
			name:   "filter by status todo",
			filter: db.IssueFilter{Statuses: []db.IssueStatus{db.IssueStatusTodo}},
			want:   []int{0, 6, 8},
		},
		{
			name:   "filter by status done",
			filter: db.IssueFilter{Statuses: []db.IssueStatus{db.IssueStatusDone}},
			want:   []int{4},
		},
		{
			name:   "filter by status killed",
			filter: db.IssueFilter{Statuses: []db.IssueStatus{db.IssueStatusKilled}},
			want:   []int{5},
		},
		{
			name:   "filter by recruiter multiple",
			filter: db.IssueFilter{RecruiterIDs: []string{"lallos", "qa"}},
			want:   []int{1, 2, 6},
		},
		{
			name:   "filter by assignee tesohh",
			filter: db.IssueFilter{AssigneeIDs: []string{"tesohh"}},
			want:   []int{0, 1, 2, 4},
		},
		{
			name:   "filter by assignee dev1",
			filter: db.IssueFilter{AssigneeIDs: []string{"dev1"}},
			want:   []int{5, 6},
		},
		{
			name:   "filter by assignee none",
			filter: db.IssueFilter{AssigneeIDs: []string{"nobody"}},
			want:   []int{},
		},
		{
			name:   "filter by tag with space trimming",
			filter: db.IssueFilter{Tags: []string{"long title"}},
			want:   []int{6},
		},
		{
			name:   "filter by category urgent",
			filter: db.IssueFilter{CategoryRoleIDs: []string{"urgent"}},
			want:   []int{7},
		},
		{
			name:   "filter by priority p0",
			filter: db.IssueFilter{PriorityRoleIDs: []string{"p0"}},
			want:   []int{6, 7},
		},
		{
			name:   "filter by title substring",
			filter: db.IssueFilter{Title: "ADS"},
			want:   []int{2, 3},
		},
		{
			name: "filter by combined status and priority",
			filter: db.IssueFilter{
				Statuses:        []db.IssueStatus{db.IssueStatusWorking},
				PriorityRoleIDs: []string{"critical"},
			},
			want: []int{9},
		},
		{
			name:   "empty filter returns all",
			filter: db.DefaultFilter(),
			want:   []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:   "no match filter",
			filter: db.IssueFilter{Title: "does not exist"},
			want:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIssues := tt.filter.Apply(testSampleIssues)
			gotIDs := []int{}
			for _, issue := range gotIssues {
				gotIDs = append(gotIDs, int(issue.ID))
			}
			if !reflect.DeepEqual(tt.want, gotIDs) {
				t.Errorf("expected IDs %v, got %v", tt.want, gotIDs)
			}
		})
	}
}
