package dataview

import (
	"crypto/rand"
	"fmt"
	"issues/v2/db"
	"issues/v2/helper"
	"issues/v2/slash"

	dg "github.com/bwmarrin/discordgo"
)

const MaxDependenciesPerPage int = 7

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
					"**Priority**: <@&%s>\n**Category**: <@&%s>",
					issue.PriorityRoleID, issue.CategoryRoleID),
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

// relationships need to have ToIssue preloaded
func MakeDependenciesContainer(issue *db.Issue, relationships db.RelationshipsByDirection) (dg.Container, bool) {
	container := dg.Container{
		AccentColor: slash.Ptr(slash.EmbedColor),
	}

	if len(relationships.Inbound) == 0 && len(relationships.Outbound) == 0 {
		return container, false
	}

	str := ""
	for _, relationship := range relationships.Inbound {
		if relationship.Kind == db.RelationshipKindDependency {

			tags := relationship.ToIssue.PrettyTags(MaxTagsCount, MaxTagLength)
			str += fmt.Sprintf("\n- %s %s %s %s",
				relationship.FromIssue.PrettyLink(len(fmt.Sprint(*relationship.FromIssue.Code))),
				relationship.FromIssue.RoleEmojis(),
				relationship.FromIssue.CutTitle(MaxTitleLength-len(tags)),
				tags,
			)

		}
	}

	// Show dependants only on the first page
	if issue.UIDepsCurrentPage == 0 && len(relationships.Inbound) > 0 {
		title := fmt.Sprintf("### Dependants `[%d]`", len(relationships.Inbound))
		container.Components = append(container.Components, dg.TextDisplay{Content: title})
		container.Components = append(container.Components, dg.TextDisplay{Content: str})
	}

	if len(relationships.Outbound) > 0 {
		completed := 0
		total := 0
		for _, r := range relationships.Outbound {
			if r.ToIssue.Status == db.IssueStatusDone {
				completed += 1
				total += 1
			} else if r.ToIssue.Status != db.IssueStatusCancelled {
				total += 1
			}
		}
		title := fmt.Sprintf("### Dependencies `[%d/%d]`", completed, total)
		container.Components = append(container.Components, dg.TextDisplay{Content: title})
	}

	paginatedDeps := helper.Paginate(relationships.Outbound, MaxDependenciesPerPage, issue.UIDepsCurrentPage)

	for _, relationship := range paginatedDeps {
		if relationship.Kind == db.RelationshipKindDependency {
			preview := ""

			switch relationship.ToIssue.Kind {
			case db.IssueKindNormal:
				tags := relationship.ToIssue.PrettyTags(MaxTagsCount, MaxTagLength)
				preview = fmt.Sprintf("- %s %s %s %s",
					relationship.ToIssue.PrettyLink(len(fmt.Sprint(*relationship.ToIssue.Code))),
					relationship.ToIssue.RoleEmojis(),
					relationship.ToIssue.CutTitle(MaxTitleLength-len(tags)),
					tags,
				)
			case db.IssueKindTask:
				preview = "- " + relationship.ToIssue.PrettyTask(MaxTitleLength+((MaxTagLength+3)*MaxTagsCount))
			}

			container.Components = append(container.Components, dg.Section{
				Components: []dg.MessageComponent{
					dg.TextDisplay{Content: preview},
				},
				Accessory: dg.Button{CustomID: "select" + rand.Text(), Label: "Select", Style: dg.SecondaryButton}, // TODO:
			})
		}
	}

	pages := helper.Pages(relationships.Outbound, MaxDependenciesPerPage)
	if pages > 1 {
		pageText := fmt.Sprintf("\n-# page %d/%d", issue.UIDepsCurrentPage+1, pages)
		container.Components = append(container.Components, dg.TextDisplay{Content: pageText})
	}

	return container, true
}

func MakeDependenciesPaginationButtons(issue *db.Issue, relationships db.RelationshipsByDirection) (dg.ActionsRow, bool) {
	if len(relationships.Outbound) < MaxDependenciesPerPage {
		return dg.ActionsRow{}, false
	}

	pages := helper.Pages(relationships.Outbound, MaxDependenciesPerPage)

	leftDisable := issue.UIDepsCurrentPage <= 0
	rightDisable := issue.UIDepsCurrentPage >= pages-1

	arrowButtons := dg.ActionsRow{
		Components: []dg.MessageComponent{
			dg.Button{Emoji: &dg.ComponentEmoji{Name: "⬅️"}, Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issue-deps-goto:%d:%d:left", issue.ID, issue.UIDepsCurrentPage-1), Disabled: leftDisable},
			dg.Button{Emoji: &dg.ComponentEmoji{Name: "➡️"}, Style: dg.SecondaryButton, CustomID: fmt.Sprintf("issue-deps-goto:%d:%d:right", issue.ID, issue.UIDepsCurrentPage+1), Disabled: rightDisable},
		},
	}

	return arrowButtons, true
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

func MakeIssueThreadDetail(issue *db.Issue, relationships db.RelationshipsByDirection, nobodyRoleID string) []dg.MessageComponent {
	allComponents := []dg.MessageComponent{
		MakeIssueMainDetail(issue, nobodyRoleID),
		dg.ActionsRow{
			Components: []dg.MessageComponent{
				makeIssueNextStateButton(issue),
				makeAssignMeButton(issue),
			},
		},
	}

	dependenciesContainer, ok := MakeDependenciesContainer(issue, relationships)
	if ok {
		allComponents = append(allComponents, dg.Separator{
			Divider: slash.Ptr(false),
			Spacing: slash.Ptr(dg.SeparatorSpacingSizeLarge),
		})
		allComponents = append(allComponents, dependenciesContainer)
	}
	dependenciesButtons, ok := MakeDependenciesPaginationButtons(issue, relationships)
	if ok {
		allComponents = append(allComponents, dependenciesButtons)
	}

	return allComponents
}
