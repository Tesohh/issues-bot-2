package handler

import (
	"fmt"
	"issues/v2/dataview"
	"issues/v2/db"
	"issues/v2/slash"

	dg "github.com/bwmarrin/discordgo"
)

func issuesShowMine(s *dg.Session, i *dg.InteractionCreate, args []string) error {
	project, err := db.Projects.
		Preload("Issues", nil).
		Preload("Issues.PriorityRole", nil).
		Preload("Issues.CategoryRole", nil).
		Where("id = ?", args[1]).
		First(db.Ctx)
	if err != nil {
		return err
	}

	issues := project.Issues
	for i := range issues {
		issues[i].Project.GuildID = project.GuildID
	}

	// create a new state
	state := db.ProjectViewState{
		ProjectID: project.ID,
		Project:   project,
		Filter: db.IssueFilter{
			Statuses:        []db.IssueStatus{db.IssueStatusTodo, db.IssueStatusWorking},
			Tags:            []string{},
			PriorityRoleIDs: []string{},
			CategoryRoleIDs: []string{},
			RecruiterIDs:    []string{i.Member.User.ID},
			AssigneeIDs:     []string{},
			Title:           "",
		},
		Sorter:      db.DefaultSorter(),
		ListNameFmt: "# My issues in %s `[%s]`",
	}
	// send the message with no buttons
	components := dataview.MakeInteractiveIssuesView(issues, &state, true)

	err = slash.ReplyWithComponents(s, i.Interaction, true, components...)
	if err != nil {
		fmt.Printf("components: %v\n", components)
		return err
	}

	msg, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		return err
	}
	state.MessageID = msg.ID
	state.ChannelID = msg.ChannelID

	state.Project = db.Project{} // im so sorry. this causes issues with gorm for some reason.
	err = db.ProjectViewStates.Create(db.Ctx, &state)
	if err != nil {
		return err
	}

	state.Project = project // im so sorry
	// make the view WITH the buttons
	components = dataview.MakeInteractiveIssuesView(issues, &state, false)
	_, err = s.InteractionResponseEdit(i.Interaction, &dg.WebhookEdit{
		Components: &components,
	})
	if err != nil {
		return err
	}
	return nil
}
