package db

import (
	"database/sql"
	"fmt"
	"issues/v2/helper"
	"strings"
	"time"
)

type IssueStatus int

const (
	IssueStatusTodo      IssueStatus = 0
	IssueStatusDoing     IssueStatus = 1
	IssueStatusDone      IssueStatus = 2
	IssueStatusCancelled IssueStatus = 3
)

var IssueStatusIcons = []string{"🟩", "🟦", "🟪", "🟥"}
var IssueStatusColors = []int{0x7cb45c, 0x54acee, 0xa98ed6, 0xdd2e44}
var IssueStatusNames = []string{"todo", "doing", "done", "cancelled"}

type IssueKind string

const (
	IssueKindNormal     IssueKind = ""
	IssueKindTask       IssueKind = "task"
	IssueKindDiscussion IssueKind = "discussion"
)

type Issue struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt sql.NullTime

	Code   *uint // for Tasks this will be nil
	Title  string
	Status IssueStatus `gorm:"check:status>=0;check:status <=3"`
	Kind   IssueKind   `gorm:"check:kind in ('', 'task',  'discussion')"`

	Tags []Tag `gorm:"many2many:issue_tags"` // comma separated

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

// func (issue *Issue) ParseTags() []string {
// 	return ParseTags(issue.Tags)
// }

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

// - [`🟩 #25`](https://example.com)
// this outputs a 84 character string, considering the code is 4 digits long
// Requires issue.Project.GuildID to be set, or else the link will be broken
func (issue *Issue) PrettyLink(longestCodeLen int) string {
	codePaddingFmt := fmt.Sprintf("%%0%dd", longestCodeLen)
	codeFmt := fmt.Sprintf(codePaddingFmt, *issue.Code)
	return fmt.Sprintf("[`%s #%s`](https://discord.com/channels/%s/%s)",
		IssueStatusIcons[issue.Status],
		codeFmt,
		issue.Project.GuildID,
		issue.ThreadID,
	)
}

func (issue *Issue) CutTitle(maxTitleLength int) string {
	return helper.StrTrunc(issue.Title, maxTitleLength)
}

func (issue *Issue) PrettyTags(maxTags int, maxTagLen int) string {
	pretties := []string{}
	for _, tag := range issue.Tags[:min(len(issue.Tags), maxTags)] {
		pretties = append(pretties, tag.Pretty(maxTagLen))
	}
	if len(issue.Tags) > maxTags {
		pretties = append(pretties, fmt.Sprintf("`[+%d]`", len(issue.Tags)-maxTags))
	}
	return strings.Join(pretties, " ")
}

// Requires PriorityRole and CategoryRole to be set
func (issue *Issue) RoleEmojis() string {
	return fmt.Sprintf("%s %s", issue.PriorityRole.Emoji, issue.CategoryRole.Emoji)
}
