package logic

import (
	"issues/v2/dataview"
	"issues/v2/db"
	"issues/v2/slash"

	dg "github.com/bwmarrin/discordgo"
)

// Requires that state.Project is set, and that Issues have issue.Project.GuildID set for proper formatting
func InitIssueView(s *dg.Session, i *dg.Interaction, state *db.ProjectViewState, ephemeral bool) error {
	project := state.Project
	// send the message with no buttons
	components := dataview.MakeInteractiveIssuesView(project.Issues, state, true)

	err := slash.ReplyWithComponents(s, i, ephemeral, components...)
	if err != nil {
		return err
	}

	msg, err := s.InteractionResponse(i)
	if err != nil {
		return err
	}
	state.MessageID = msg.ID
	state.ChannelID = msg.ChannelID

	state.Project = db.Project{} // im so sorry. this causes issues with gorm for some reason.
	err = db.ProjectViewStates.Create(db.Ctx, state)
	if err != nil {
		return err
	}

	state.Project = project // im so sorry

	// make the view WITH the buttons
	components = dataview.MakeInteractiveIssuesView(project.Issues, state, false)
	_, err = s.InteractionResponseEdit(i, &dg.WebhookEdit{
		Components: &components,
	})
	if err != nil {
		return err
	}
	return nil
}

// Requires that state.Project is set, and that Issues have issue.Project.GuildID set for proper formatting
// does the same thing as InitIssueView, but is not attached to a reply
func InitIssueViewDetached(s *dg.Session, i *dg.Interaction, channelID string, state *db.ProjectViewState) error {
	project := state.Project
	// send the message with no buttons
	components := dataview.MakeInteractiveIssuesView(project.Issues, state, true)

	msg, err := s.ChannelMessageSendComplex(channelID, &dg.MessageSend{
		Components: components,
		Flags:      dg.MessageFlagsIsComponentsV2,
	})
	if err != nil {
		return err
	}

	state.MessageID = msg.ID
	state.ChannelID = msg.ChannelID

	state.Project = db.Project{} // im so sorry. this causes issues with gorm for some reason.
	err = db.ProjectViewStates.Create(db.Ctx, state)
	if err != nil {
		return err
	}

	state.Project = project // im so sorry

	// make the view WITH the buttons
	components = dataview.MakeInteractiveIssuesView(project.Issues, state, false)
	_, err = s.ChannelMessageEditComplex(&dg.MessageEdit{
		Components: &components,
		ID:         msg.ID,
		Channel:    channelID,
	})
	if err != nil {
		return err
	}
	return nil
}
