package handler

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/logic"

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
	fmt.Printf("assigneeIDs: %v, botid: %v\n", assigneeIDs, s.State.User.ID)
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
	return nil
}
