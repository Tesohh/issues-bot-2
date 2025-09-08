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

	return dg.Container{
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
}

func MakeIssueThreadDetail(issue *db.Issue, nobodyRoleID string) []dg.MessageComponent {
	allComponents := []dg.MessageComponent{
		// ...
		MakeIssueMainDetail(issue, nobodyRoleID),
		// ...
	}

	return allComponents
}
