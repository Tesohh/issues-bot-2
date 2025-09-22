package command

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/logic"
	"issues/v2/slash"

	dg "github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var Task = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "task",
		Description: "task adding and management commands",
		Options: []*dg.ApplicationCommandOption{
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "new",
				Description: "adds a new task under the current issue",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "title",
						Description: "the title of the task",
						Required:    true,
					},
				},
			},
		},
	},
	Func: func(s *dg.Session, i *dg.Interaction) error {
		subcommand := i.ApplicationCommandData().Options[0]
		opts := slash.GetOptionMapRaw(subcommand.Options)

		issue, err := db.Issues.
			Preload("Tags", nil).
			Preload("AssigneeUsers", nil).
			Preload("Project", func(db gorm.PreloadBuilder) error {
				db.Select("ID", "Prefix", "GuildID")
				return nil
			}).
			Where("thread_id = ?", i.ChannelID).
			First(db.Ctx)

		if err == gorm.ErrRecordNotFound {
			return ErrNotInIssueThread
		} else if err != nil {
			return err
		}

		switch subcommand.Name {
		case "new":
			title := opts["title"].StringValue()
			err = TaskNew(s, i, &issue, title)
		}

		if err != nil {
			return err
		}

		err = logic.UpdateEverythingAboutSingleIssue(s, i.GuildID, &issue)
		if err != nil {
			return err
		}

		return nil
	},
}

func TaskNew(s *dg.Session, i *dg.Interaction, issue *db.Issue, title string) error {
	// create the task in the DB
	task := db.Issue{
		Code:            nil,
		Title:           title,
		Status:          db.IssueStatusTodo,
		Kind:            db.IssueKindTask,
		ProjectID:       issue.ProjectID,
		RecruiterUserID: issue.RecruiterUserID,
		AssigneeUsers:   issue.AssigneeUsers,
		CategoryRoleID:  issue.CategoryRoleID,
		PriorityRoleID:  issue.PriorityRoleID,
	}

	err := db.Issues.Create(db.Ctx, &task)
	if err != nil {
		return err
	}

	// create the relationship in the DB
	relationship := db.Relationship{
		FromIssueID: issue.ID,
		ToIssueID:   task.ID,
		Kind:        db.RelationshipKindDependency,
	}

	err = db.Relationships.Create(db.Ctx, &relationship)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("<@%s> added task `%s`", i.Member.User.ID, task.CutTitle(25))
	return slash.ReplyWithText(s, i, msg, false)
}
