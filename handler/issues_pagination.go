package handler

import (
	"issues/v2/dataview"
	"issues/v2/db"
	"issues/v2/helper"
	"issues/v2/logic"

	dg "github.com/bwmarrin/discordgo"
)

func issuesBigLeft(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	_, err := db.ProjectViewStates.Where("message_id = ?", args[1]).Update(db.Ctx, "current_page", 0)
	if err != nil {
		return err
	}

	return logic.UpdateInteractiveIssuesView(s, args[1], false)
}
func issuesLeft(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	state, err := db.ProjectViewStates.
		Select("current_page").
		Where("message_id = ?", args[1]).
		First(db.Ctx)
	if err != nil {
		return err
	}
	_, err = db.ProjectViewStates.Where("message_id = ?", args[1]).Update(db.Ctx, "current_page", state.CurrentPage-1)
	if err != nil {
		return err
	}

	return logic.UpdateInteractiveIssuesView(s, args[1], false)
}
func issuesRight(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	state, err := db.ProjectViewStates.
		Select("current_page").
		Where("message_id = ?", args[1]).
		First(db.Ctx)
	if err != nil {
		return err
	}
	_, err = db.ProjectViewStates.Where("message_id = ?", args[1]).Update(db.Ctx, "current_page", state.CurrentPage+1)
	if err != nil {
		return err
	}

	return logic.UpdateInteractiveIssuesView(s, args[1], false)
}
func issuesBigRight(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	// PERF: iper smarcio
	// we SADLY need to get the number of pages by using the filter etc

	state, err := db.ProjectViewStates.
		Select("filter").
		Preload("Project.Issues", nil).
		Where("message_id = ?", args[1]).
		First(db.Ctx)
	if err != nil {
		return err
	}
	filteredIssues := state.Filter.Apply(state.Project.Issues)
	pages := helper.Pages(filteredIssues, dataview.MaxIssuesPerPage)

	_, err = db.ProjectViewStates.Where("message_id = ?", args[1]).Update(db.Ctx, "current_page", pages)
	if err != nil {
		return err
	}

	return logic.UpdateInteractiveIssuesView(s, args[1], false)
}
