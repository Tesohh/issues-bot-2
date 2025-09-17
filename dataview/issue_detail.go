package dataview

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/slash"

	dg "github.com/bwmarrin/discordgo"
)

func MakeIssueMainDetail(issue *db.Issue, nobodyRoleID string) dg.Container {
	assigneeIDs := []string{}
	for _, user := range issue.AssigneeUsers {
		assigneeIDs = append(assigneeIDs, user.ID)
	}

	assigneeStr := ""
	if len(assigneeIDs) == 0 {
		assigneeStr = fmt.Sprintf("<@&%s>", nobodyRoleID)
	} else {
		assigneeStr = slash.MentionMany(assigneeIDs, "@", ", ")
	}

	container := dg.Container{
		AccentColor: slash.Ptr(slash.EmbedColor),
		Components: []dg.MessageComponent{
			dg.TextDisplay{
				Content: fmt.Sprintf("## `%s` %s %s", issue.HumanCode(), db.IssueStatusIcons[issue.Status], issue.Title),
			},

			dg.Separator{},

			dg.TextDisplay{
				Content: fmt.Sprintf(
					"**Category**: <@&%s>\n**Priority**: <@&%s>",
					issue.CategoryRoleID, issue.PriorityRoleID),
			},

			dg.Separator{},

			dg.TextDisplay{
				Content: fmt.Sprintf("**Recruiter**: <@%s>\n**Assignee(s)**: %s", issue.RecruiterUserID, assigneeStr),
			},
		},
	}

	if len(issue.Tags) > 0 {
		container.Components = append(container.Components, dg.Separator{}, dg.TextDisplay{
			Content: fmt.Sprintf("**Tags**: %s", issue.PrettyTags(999, 999)),
		})
	}

	return container
}

func makeIssueNextStateButton(issue *db.Issue) dg.Button {
	label := ""
	style := dg.SecondaryButton
	status := db.IssueStatusTodo
	disabled := false
	emoji := ""

	switch issue.Status {
	case db.IssueStatusTodo:
		label = "Mark as Doing"
		// style = dg.PrimaryButton
		status = db.IssueStatusDoing
		emoji = db.IssueStatusIcons[db.IssueStatusDoing]
	case db.IssueStatusDoing:
		label = "Mark as Done"
		// style = dg.SuccessButton
		status = db.IssueStatusDone
		emoji = db.IssueStatusIcons[db.IssueStatusDone]
	case db.IssueStatusDone:
		label = "Revert to Todo"
		// style = dg.SecondaryButton
		status = db.IssueStatusTodo
		emoji = db.IssueStatusIcons[db.IssueStatusTodo]
	case db.IssueStatusCancelled:
		label = "Mark as Cancelled"
		// style = dg.DangerButton
		status = db.IssueStatusCancelled
		emoji = db.IssueStatusIcons[db.IssueStatusCancelled]
		disabled = true
	}

	return dg.Button{
		Label:    label,
		Style:    style,
		CustomID: fmt.Sprintf("issue-set-status:%d:%d", issue.ID, status),
		Disabled: disabled,
		Emoji:    &dg.ComponentEmoji{Name: emoji},
	}
}

func makeAssignMeButton(issue *db.Issue) dg.Button {
	return dg.Button{
		Label:    "Assign me",
		Style:    dg.PrimaryButton,
		CustomID: fmt.Sprintf("issue-toggle-author-assignee:%d", issue.ID),
		Disabled: false,
	}
}

func MakeIssueThreadDetail(issue *db.Issue, nobodyRoleID string) []dg.MessageComponent {
	allComponents := []dg.MessageComponent{
		MakeIssueMainDetail(issue, nobodyRoleID),
		dg.ActionsRow{
			Components: []dg.MessageComponent{
				makeIssueNextStateButton(issue),
				makeAssignMeButton(issue),
			},
		},
	}

	return allComponents
}
