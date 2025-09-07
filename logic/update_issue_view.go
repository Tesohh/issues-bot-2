package logic

import (
	"issues/v2/dataview"
	"issues/v2/db"

	dg "github.com/bwmarrin/discordgo"
)

// updates a single view
// use this when changing settings of a view for example
// also sets the page to 0 if page0 is true
func UpdateInteractiveIssuesView(s *dg.Session, messageID string, page0 bool) error {
	if page0 {
		_, err := db.ProjectViewStates.Where("message_id = ?", messageID).Update(db.Ctx, "current_page", 0)
		if err != nil {
			return err
		}
	}

	state, err := db.ProjectViewStates.
		Preload("Project.Issues", nil).
		Preload("Project.Issues.PriorityRole", nil).
		Preload("Project.Issues.CategoryRole", nil).
		Where("message_id = ?", messageID).
		First(db.Ctx)
	if err != nil {
		return err
	}

	components := dataview.MakeInteractiveIssuesView(state.Project.Issues, &state, false)
	_, err = s.ChannelMessageEditComplex(&dg.MessageEdit{
		Components: &components,
		ID:         messageID,
		Channel:    state.ChannelID,
	})
	return err
}

// updates all views linked to a project
// use this when issues in a project are updated
// does not change the page, as someone updating an issue would change the page and would be annoying
func UpdateAllInteractiveIssuesViews(s *dg.Session, projectID uint) error {
	project, err := db.Projects.Preload("Issues", nil).
		Preload("Issues.PriorityRole", nil).
		Preload("Issues.CategoryRole", nil).
		Where("id = ?", projectID).
		First(db.Ctx)
	if err != nil {
		return err
	}

	states, err := db.ProjectViewStates.Where("project_id = ?", projectID).Find(db.Ctx)
	if err != nil {
		return err
	}
	// TODO: add warning message for having too many lsits as this can be slow
	for _, state := range states {
		state.Project = project
		components := dataview.MakeInteractiveIssuesView(project.Issues, &state, false)
		_, err = s.ChannelMessageEditComplex(&dg.MessageEdit{
			Components: &components,
			ID:         state.MessageID,
			Channel:    state.ChannelID,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
