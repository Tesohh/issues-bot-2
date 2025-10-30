package handler

import (
	"issues/v2/db"
	"issues/v2/logic"
	"log/slog"
	"strconv"

	dg "github.com/bwmarrin/discordgo"
)

func issuesFilterPeopleSubmit(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	// extract user IDs
	data := i.ModalSubmitData()
	recruiterIDs := data.Components[0].(*dg.Label).Component.(*dg.SelectMenu).Values
	assigneeIDs := data.Components[1].(*dg.Label).Component.(*dg.SelectMenu).Values

	// update filter to use these recruiter and assignee IDs
	state, err := db.ProjectViewStates.Where("message_id = ?", args[1]).First(db.Ctx)
	if err != nil {
		return err
	}
	state.Filter.RecruiterIDs = recruiterIDs
	state.Filter.AssigneeIDs = assigneeIDs
	if len(assigneeIDs) == 1 && assigneeIDs[0] == s.State.User.ID {
		state.Filter.Nobody = true
		state.Filter.AssigneeIDs = []string{}
	} else {
		state.Filter.Nobody = false
		state.Filter.AssigneeIDs = assigneeIDs
	}
	_, err = db.ProjectViewStates.Updates(db.Ctx, state)
	if err != nil {
		return err
	}

	return logic.UpdateInteractiveIssuesView(s, args[1], true)
}

func issuesFilterDataSubmit(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	// extract data
	data := i.ModalSubmitData()
	var (
		priorities  = data.Components[0].(*dg.Label).Component.(*dg.SelectMenu).Values
		categories  = data.Components[1].(*dg.Label).Component.(*dg.SelectMenu).Values
		statusesRaw = data.Components[2].(*dg.Label).Component.(*dg.SelectMenu).Values
		tagsRaw     = data.Components[3].(*dg.Label).Component.(*dg.TextInput).Value
		title       = data.Components[4].(*dg.Label).Component.(*dg.TextInput).Value
	)

	statuses := []db.IssueStatus{}
	for _, s := range statusesRaw {
		i, err := strconv.Atoi(s)
		if err != nil {
			slog.Warn("strconv.Atoi error in issuesFilterDataSubmit", "err", err, "i", i)
			continue
		}
		statuses = append(statuses, db.IssueStatus(i))
	}

	// update filter
	state, err := db.ProjectViewStates.Where("message_id = ?", args[1]).First(db.Ctx)
	if err != nil {
		return err
	}

	state.Filter.PriorityRoleIDs = priorities
	state.Filter.CategoryRoleIDs = categories
	state.Filter.Statuses = statuses
	state.Filter.Tags = db.ParseTags(tagsRaw)
	state.Filter.Title = title

	_, err = db.ProjectViewStates.Updates(db.Ctx, state)
	if err != nil {
		return err
	}

	return logic.UpdateInteractiveIssuesView(s, args[1], true)
}
