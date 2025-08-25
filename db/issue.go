package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type IssueStatus uint8

const (
	IssueStatusTodo    IssueStatus = 0
	IssueStatusWorking IssueStatus = 1
	IssueStatusDone    IssueStatus = 2
	IssueStatusKilled  IssueStatus = 3
)

var IssueStatusIcons = [4]string{"🟩", "🟦", "🟪", "🟥"}
var IssueStatusColors = [4]int{0x7cb45c, 0x54acee, 0xa98ed6, 0xdd2e44}
var IssueStatusNames = [4]string{"todo", "working", "done", "killed"}

type IssueKind string

const (
	IssueKindNormal     IssueKind = ""
	IssueKindTask       IssueKind = "task"
	IssueKindDiscussion IssueKind = "discussion"
)

// (Id, Code, Title, Status, Tags, Kind {Normal, Task, Discussion},
// ProjectID, RecruiterID, AssigneeIDs, ThreadID, MessageID, CategoryRoleID, PriorityRoleID)
// - [ ] id = auto increment
// - [ ] Code = max(Code)+1 for issues and discussions but NULL for tasks

type Issue struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	CompletedAt sql.NullTime

	Code   *uint // for Tasks this will be nil
	Title  string
	Tags   string      // comma separated
	Status IssueStatus `gorm:"check:status>=0;check:status <=3"`
	Kind   IssueKind   `gorm:"check:kind in ('', 'task',  'discussion')"`

	ProjectID uint
	Project   Project

	RecruiterUserID string
	RecruiterUser   User
	AssigneeUsers   []User `gorm:"many2many:issue_assignees;"`

	CategoryRoleID string
	CategoryRole   Role

	PriorityRoleID string
	PriorityRole   Role

	ThreadID  string
	MessageID string
}

// Requires issue.Project.Prefix to be set, or else the prefix will be ???
func (issue *Issue) HumanCode() string {
	projectName := "???"
	if len(issue.Project.Prefix) > 0 {
		projectName = strings.ToUpper(issue.Project.Prefix)
	}
	return fmt.Sprintf("#%s-%d", projectName, *issue.Code)
}

// Requires issue.Project.Prefix to be set, or else the prefix will be ???
func (issue *Issue) ChannelName() string {
	return fmt.Sprintf("%s %s %s", issue.HumanCode(), IssueStatusIcons[issue.Status], issue.Title)
}

// - [`🟩 #25`](https://example.com) add Gòrni
// Requires issue.Project.GuildID to be set, or else the link will be broken
func (issue *Issue) PrettyLink() string {
	return fmt.Sprintf("[`%s #%d`](https://discord.com/channels/%s/%s) %s", IssueStatusIcons[issue.Status], *issue.Code, issue.Project.GuildID, issue.ThreadID, issue.Title)
}
