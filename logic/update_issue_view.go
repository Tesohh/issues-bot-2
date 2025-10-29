package logic

import (
	"issues/v2/dataview"
	"issues/v2/db"
	"log/slog"
	"time"

	dg "github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
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
		Preload("Project.Issues", func(query gorm.PreloadBuilder) error {
			query.Where("kind = ?", db.IssueKindNormal)
			return nil
		}).
		Preload("Project.Issues.Tags", nil).
		Preload("Project.Issues.AssigneeUsers", nil).
		Preload("Project.Issues.PriorityRole", nil).
		Preload("Project.Issues.CategoryRole", nil).
		Where("message_id = ?", messageID).
		First(db.Ctx)
	if err != nil {
		return err
	}

	for i := range state.Project.Issues {
		state.Project.Issues[i].Project.GuildID = state.Project.GuildID
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
// also purges old lists
func UpdateAllInteractiveIssuesViews(s *dg.Session, projectID uint) error {
	project, err := db.Projects.
		Preload("Issues", func(query gorm.PreloadBuilder) error {
			query.Where("kind = ?", db.IssueKindNormal)
			return nil
		}).
		Preload("Issues.Tags", nil).
		Preload("Issues.PriorityRole", nil).
		Preload("Issues.CategoryRole", nil).
		Where("id = ?", projectID).
		First(db.Ctx)
	if err != nil {
		return err
	}
	for i := range project.Issues {
		project.Issues[i].Project.GuildID = project.GuildID
	}

	states, err := db.ProjectViewStates.Where("project_id = ?", projectID).Find(db.Ctx)
	if err != nil {
		return err
	}

	for _, state := range states {
		if !state.Permanent && time.Since(state.UpdatedAt) > 24*time.Hour {
			state.DeletedAt.Valid = true

			_, err = db.ProjectViewStates.Where("message_id = ?", state.MessageID).Delete(db.Ctx)
			if err != nil {
				slog.Error("error while deleting state", "err", err, "state", state)
			}
		}

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
