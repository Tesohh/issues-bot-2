package command

import (
	"issues/v2/db"
	"issues/v2/slash"
	"slices"

	"github.com/bwmarrin/discordgo"
	dg "github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var codeOpt = dg.ApplicationCommandOption{
	Type:        dg.ApplicationCommandOptionInteger,
	Name:        "code",
	Description: "the code of the issue to edit. is inferred if you're in an issue thread",
	Required:    false,
	MinValue:    slash.Ptr(float64(0)),
	MaxValue:    0,
}

var Issue = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "issue",
		Description: "various issue editing commands",
		Options: []*dg.ApplicationCommandOption{
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "assign",
				Description: "assigns the issue to a user and removes them if they are already assigned",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionUser,
						Name:        "assignee",
						Description: "user to assign/remove",
						Required:    true,
					},
					&codeOpt,
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "category",
				Description: "changes the category of the issue",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "role to set as category",
						Required:    true,
					},
					&codeOpt,
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "priority",
				Description: "changes the priority of the issue",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "role to set as priority",
						Required:    true,
					},
					&codeOpt,
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "rename",
				Description: "changes the title of the issue",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "title",
						Description: "the new title",
						Required:    true,
					},
					&codeOpt,
				},
			},
		},
	},
	Func: func(s *dg.Session, i *dg.Interaction) error {
		subcommand := i.ApplicationCommandData().Options[0]
		options := slash.GetOptionMapRaw(subcommand.Options)

		var query gorm.ChainInterface[db.Issue]
		if codeOpt, ok := options["code"]; ok {
			code := codeOpt.IntValue()
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				return err
			}

			project, err := db.Projects.
				Select("id").
				Where("discord_category_channel_id = ?", channel.ParentID).
				Or("issues_input_channel_id = ?", channel.ParentID).
				First(db.Ctx)
			if err != nil {
				return err
			}

			query = db.Issues.
				Where("code = ?", code).
				Where("project_id = ?", project.ID)
		} else {
			query = db.Issues.Where("thread_id = ?", i.ChannelID)
		}
		issue, err := query.First(db.Ctx)
		if err == gorm.ErrRecordNotFound {
			return ErrNotInIssueThread
		}

		switch subcommand.Name {
		case "assign":
			assignee := options["assignee"].UserValue(s)
			err = IssueAssign(s, i, &issue, assignee)
		}

		if err != nil {
			return err
		}

		// also update the issue view if there were no errors
		return nil
	},
}

func IssueAssign(s *discordgo.Session, i *discordgo.Interaction, issue *db.Issue, assignee *dg.User) error {
	index := slices.IndexFunc(issue.AssigneeUsers, func(user db.User) bool {
		return user.ID == assignee.ID
	})

	if index == -1 {
		issue.AssigneeUsers = append(issue.AssigneeUsers, db.User{ID: assignee.ID})
	} else {
		issue.AssigneeUsers = slices.Delete(issue.AssigneeUsers, index, index+1)
	}

	_, err := db.Issues.Updates(db.Ctx, *issue)
	if err != nil {
		return err
	}

	return err
}
