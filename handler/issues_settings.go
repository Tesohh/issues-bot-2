package handler

import (
	"issues/v2/data"
	"issues/v2/db"
	"issues/v2/logic"
	"slices"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

func issuesGoto(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	_, err := db.ProjectViewStates.Where("message_id = ?", args[1]).Update(db.Ctx, "current_page", args[2])
	if err != nil {
		return err
	}

	return logic.UpdateInteractiveIssuesView(s, args[1], false)
}

func issuesSetStatuses(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	state, err := db.ProjectViewStates.Where("message_id = ?", args[1]).First(db.Ctx)
	if err != nil {
		return err
	}
	state.Filter.Statuses = []db.IssueStatus{}
	for name := range strings.SplitSeq(args[2], ",") {
		state.Filter.Statuses = append(state.Filter.Statuses, db.IssueStatus(slices.Index(db.IssueStatusNames, name)))
	}

	_, err = db.ProjectViewStates.Updates(db.Ctx, state)
	if err != nil {
		return err
	}

	return logic.UpdateInteractiveIssuesView(s, args[1], true)
}

func issuesSortBy(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	state, err := db.ProjectViewStates.Where("message_id = ?", args[1]).First(db.Ctx)
	if err != nil {
		return err
	}

	state.Sorter.SortBy = db.IssueSortBy(args[2])

	_, err = db.ProjectViewStates.Updates(db.Ctx, state)
	if err != nil {
		return err
	}

	return logic.UpdateInteractiveIssuesView(s, args[1], true)
}

func issuesOrder(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	state, err := db.ProjectViewStates.Where("message_id = ?", args[1]).First(db.Ctx)
	if err != nil {
		return err
	}

	state.Sorter.SortOrder = db.SortOrder(args[2])

	_, err = db.ProjectViewStates.Updates(db.Ctx, state)
	if err != nil {
		return err
	}

	return logic.UpdateInteractiveIssuesView(s, args[1], true)
}

func issuesFilterPeople(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	state, err := db.ProjectViewStates.Where("message_id = ?", args[1]).First(db.Ctx)
	if err != nil {
		return err
	}

	recruitersDefault := []dg.SelectMenuDefaultValue{}
	for _, id := range state.Filter.RecruiterIDs {
		recruitersDefault = append(recruitersDefault, dg.SelectMenuDefaultValue{ID: id, Type: dg.SelectMenuDefaultValueUser})
	}

	assigneesDefault := []dg.SelectMenuDefaultValue{}
	if state.Filter.Nobody {
		assigneesDefault = append(assigneesDefault, dg.SelectMenuDefaultValue{ID: s.State.User.ID, Type: dg.SelectMenuDefaultValueUser})
	} else {
		for _, id := range state.Filter.AssigneeIDs {
			assigneesDefault = append(assigneesDefault, dg.SelectMenuDefaultValue{ID: id, Type: dg.SelectMenuDefaultValueUser})
		}
	}

	err = s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
		Type: dg.InteractionResponseModal,
		Data: &dg.InteractionResponseData{
			CustomID: "issues_filter_people_submit:" + i.Message.ID,
			Title:    "Filter people",
			Flags:    dg.MessageFlagsIsComponentsV2,
			Components: []dg.MessageComponent{
				dg.Label{
					Label:       "Recruiters",
					Description: "People who created the issue",
					Component: dg.SelectMenu{
						MenuType:      dg.UserSelectMenu,
						CustomID:      "recruiters",
						Placeholder:   "filter by recruiters...",
						MaxValues:     10,
						Required:      false,
						DefaultValues: recruitersDefault,
					},
				},
				dg.Label{
					Label:       "Assignees",
					Description: "@YIELD for issues with no assignees (sorry for jank)",
					Component: dg.SelectMenu{
						MenuType:      dg.UserSelectMenu,
						CustomID:      "assignees",
						Placeholder:   "filter by assignees...",
						MaxValues:     10,
						Required:      false,
						DefaultValues: assigneesDefault,
					},
				},
			},
		},
	})
	return err
}

func issuesFilterData(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	state, err := db.ProjectViewStates.Where("message_id = ?", args[1]).First(db.Ctx)
	if err != nil {
		return err
	}

	priorityRoles, err := db.Roles.Where("guild_id = ?", i.GuildID).Where("kind = ?", db.RoleKindPriority).Find(db.Ctx)
	if err != nil {
		return err
	}
	priorityChoices := []dg.SelectMenuOption{}
	for _, r := range priorityRoles {
		priorityChoices = append(priorityChoices,
			dg.SelectMenuOption{
				Label:   strings.ToUpper(r.Key),
				Emoji:   &dg.ComponentEmoji{Name: r.Emoji},
				Value:   r.Key,
				Default: slices.Contains(state.Filter.PriorityRoleIDs, r.ID),
			},
		)
	}

	categoryRoles, err := db.Roles.Where("guild_id = ?", i.GuildID).Where("kind = ?", db.RoleKindCategory).Find(db.Ctx)
	if err != nil {
		return err
	}
	categoryChoices := []dg.SelectMenuOption{}
	for _, r := range categoryRoles {
		categoryChoices = append(categoryChoices,
			dg.SelectMenuOption{
				Label:   strings.ToUpper(r.Key),
				Emoji:   &dg.ComponentEmoji{Name: r.Emoji},
				Value:   r.Key,
				Default: slices.Contains(state.Filter.CategoryRoleIDs, r.ID),
			},
		)
	}

	statusChoices := []dg.SelectMenuOption{}
	for i, choice := range data.StatusOptionSelectChoices {
		choice.Default = slices.Contains(state.Filter.Statuses, db.IssueStatus(i))
		statusChoices = append(statusChoices, choice)
	}

	defaultTags := strings.Join(state.Filter.Tags, ", ")

	defaultTitle := state.Filter.Title

	err = s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
		Type: dg.InteractionResponseModal,
		Data: &dg.InteractionResponseData{
			CustomID: "issues_filter_data_submit:" + i.Message.ID,
			Title:    "Filter data",
			Flags:    dg.MessageFlagsIsComponentsV2,
			Components: []dg.MessageComponent{
				dg.Label{
					Label: "Priority",
					Component: dg.SelectMenu{
						MenuType:    dg.StringSelectMenu,
						CustomID:    "priorities",
						Placeholder: "filter by priority...",
						MaxValues:   4,
						Required:    false,
						Options:     priorityChoices,
					},
				},
				dg.Label{
					Label: "Category",
					Component: dg.SelectMenu{
						MenuType:    dg.StringSelectMenu,
						CustomID:    "categories",
						Placeholder: "filter by category...",
						MaxValues:   4,
						Required:    false,
						Options:     categoryChoices,
					},
				},
				dg.Label{
					Label: "Status",
					Component: dg.SelectMenu{
						MenuType:    dg.StringSelectMenu,
						CustomID:    "statuses",
						Placeholder: "filter by status...",
						MaxValues:   4,
						Required:    false,
						Options:     statusChoices,
					},
				},
				dg.Label{
					Label: "Tags",
					Component: dg.TextInput{
						CustomID:    "tags",
						Placeholder: "tag1, tag2, tag3...",
						Style:       dg.TextInputShort,
						Required:    false,
						Value:       defaultTags,
					},
				},
				dg.Label{
					Label: "title",
					Component: dg.TextInput{
						CustomID:    "title",
						Placeholder: "filter by title...",
						Style:       dg.TextInputShort,
						Required:    false,
						Value:       defaultTitle,
					},
				},
			},
		},
	})
	return err
}
