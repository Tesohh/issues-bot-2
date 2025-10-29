package handler

import (
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
	// err := s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
	// 	Type: dg.InteractionResponseModal,
	// 	Data: &dg.InteractionResponseData{
	// 		CustomID: "issues_filter_people_submit:" + i.Message.ChannelID,
	// 		Title:    "Filter people",
	// 		Flags:    dg.MessageFlagsIsComponentsV2,
	// 		Components: []dg.MessageComponent{
	// 			dg.Label{
	// 				Label: "Recruiters",
	// 				Component: dg.SelectMenu{
	// 					MinValues:   slash.Ptr(0),
	// 					MenuType:    dg.UserSelectMenu,
	// 					CustomID:    "recruiters",
	// 					Placeholder: "filter by recruiters...",
	// 					MaxValues:   3,
	// 					Required:    false,
	// 				},
	// 			},
	// 		},
	// 	},
	// })
	// return err
	return nil
}

func issuesFilterData(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	return nil
}
