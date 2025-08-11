package db

import (
	"time"
)

type IssueStatus uint8

const (
	IssueStatusTodo     IssueStatus = 0
	IssueStatusDoing    IssueStatus = 1
	IssueStatusDone     IssueStatus = 2
	IssueStatusCanceled IssueStatus = 3
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
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time

	Code   *uint // for Tasks this will be nil
	Title  string
	Tags   string // comma separated
	Status IssueStatus
	Kind   IssueKind

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
