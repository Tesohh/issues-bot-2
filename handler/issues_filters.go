package handler

import (
	"issues/v2/db"
	"issues/v2/logic"
	"slices"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

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

func issuesFilters(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	return nil
}

func issuesShowMine(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	return nil
}
