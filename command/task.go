package command

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/logic"
	"issues/v2/slash"
	"log/slog"

	dg "github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var taskOpt = dg.ApplicationCommandOption{
	Type:         dg.ApplicationCommandOptionString,
	Name:         "task",
	Description:  "which task to act upon",
	Required:     true,
	Autocomplete: true,
}

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
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "promote",
				Description: "promotes a task into an issue",
				Options:     []*dg.ApplicationCommandOption{&taskOpt},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "toggle",
				Description: "toggles task between todo and done",
				Options:     []*dg.ApplicationCommandOption{&taskOpt},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "rename",
				Description: "changes the title of the task",
				Options: []*dg.ApplicationCommandOption{
					&taskOpt,
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "title",
						Description: "the new title to set",
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
				db.Select("ID", "Prefix", "GuildID", "IssuesInputChannelID")
				return nil
			}).
			Where("thread_id = ?", i.ChannelID).
			First(db.Ctx)

		if err == gorm.ErrRecordNotFound {
			return ErrNotInIssueThread
		} else if err != nil {
			return err
		}

		task := db.Issue{}
		if subcommand.Name != "new" {
			id := opts["task"].StringValue()
			task, err = db.Issues.Where("id = ?", id).First(db.Ctx)
			if err != nil {
				return err
			}

			task.Project.IssuesInputChannelID = issue.Project.IssuesInputChannelID
			task.Project.Prefix = issue.Project.Prefix
		}

		switch subcommand.Name {
		case "new":
			title := opts["title"].StringValue()
			err = TaskNew(s, i, &issue, title)
		case "promote":
			err = TaskPromote(s, i, &task)
		case "toggle":
			err = TaskToggle(s, i, &task)
		case "rename":
			title := opts["title"].StringValue()
			err = TaskRename(s, i, &task, title)
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

func TaskPromote(s *dg.Session, i *dg.Interaction, task *db.Issue) error {
	task.Kind = db.IssueKindNormal
	_, err := db.Issues.Where("id = ?", task.ID).Update(db.Ctx, "kind", db.IssueKindNormal)
	if err != nil {
		return err
	}

	code, err := logic.GetIssueCode(task)
	if err != nil {
		return fmt.Errorf("error in issue db insertion: %w", err)
	}
	task.Code = &code
	_, err = db.Issues.Where("id = ?", task.ID).Update(db.Ctx, "code", code)
	if err != nil {
		return err
	}

	thread, err := s.ThreadStart(task.Project.IssuesInputChannelID, task.ChannelName(), dg.ChannelTypeGuildPublicThread, 10080)
	if err != nil {
		return err
	}
	task.ThreadID = thread.ID
	_, err = db.Issues.Where("id = ?", task.ID).Update(db.Ctx, "thread_id", thread.ID)
	if err != nil {
		return err
	}

	err = s.ChannelMessageDelete(task.Project.IssuesInputChannelID, thread.ID)
	if err != nil {
		slog.Warn("couldn't delete thread start message. no big deal", "err", err)
	}

	guild, err := db.Guilds.Select("nobody_role_id").Where("id = ?", i.GuildID).First(db.Ctx)
	if err != nil {
		return err
	}

	err = logic.InitIssueThread(task, db.RelationshipsByDirection{}, &guild, thread, s)
	if err != nil {
		return err
	}
	_, err = db.Issues.Where("id = ?", task.ID).Update(db.Ctx, "message_id", task.MessageID)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("<@%s> promoted task `%s` to <#%s>", i.Member.User.ID, task.CutTitle(25), task.ThreadID)
	return slash.ReplyWithText(s, i, msg, false)
}

func TaskToggle(s *dg.Session, i *dg.Interaction, task *db.Issue) error {
	var status db.IssueStatus
	var word string

	switch task.Status {
	case db.IssueStatusTodo:
		status = db.IssueStatusDone
		word = "checked"
	case db.IssueStatusDone:
		status = db.IssueStatusTodo
		word = "unchecked"
	}
	_, err := db.Issues.Where("id = ?", task.ID).Update(db.Ctx, "status", status)
	if err != nil {
		return err
	}

	task.Status = status

	msg := fmt.Sprintf("<@%s> %s task `%s`", i.Member.User.ID, word, task.CutTitle(25))
	return slash.ReplyWithText(s, i, msg, false)
}

func TaskRename(s *dg.Session, i *dg.Interaction, task *db.Issue, title string) error {
	oldTask := *task
	_, err := db.Issues.Where("id = ?", task.ID).Update(db.Ctx, "title", title)
	if err != nil {
		return err
	}

	task.Title = title

	msg := fmt.Sprintf("<@%s> renamed task `%s` to `%s`", i.Member.User.ID, oldTask.CutTitle(25), task.CutTitle(25))
	return slash.ReplyWithText(s, i, msg, false)
}
