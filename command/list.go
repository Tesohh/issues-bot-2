package command

import (
	"issues/v2/db"
	"issues/v2/logic"
	"issues/v2/slash"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var List = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "list",
		Description: "lists various things",
		Options: []*dg.ApplicationCommandOption{
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "issues",
				Description: "lists issues in a project",
				Options: []*dg.ApplicationCommandOption{
					{
						Type:        dg.ApplicationCommandOptionString,
						Name:        "prefix",
						Description: "the project's prefix",
						MinLength:   slash.Ptr(3),
						MaxLength:   3,
					},
				},
			},
		},
	},
	Disabled: false,
	Func: func(s *dg.Session, i *dg.Interaction) error {
		subcommand := i.ApplicationCommandData().Options[0]
		options := slash.GetOptionMapRaw(subcommand.Options)

		switch subcommand.Name {
		case "issues":
			// get the project from the channel
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				return err
			}

			query := db.Projects.
				Preload("Issues", nil).
				Preload("Issues.PriorityRole", nil).
				Preload("Issues.CategoryRole", nil)

			project, err := query.Where("discord_category_channel_id = ?", channel.ParentID).First(db.Ctx)
			if err == gorm.ErrRecordNotFound {
				// if that didn't work, get from the prefix option
				prefix := ""
				if prefixOpt, ok := options["prefix"]; ok {
					prefix = strings.ToLower(prefixOpt.StringValue())
				} else {
					return ErrPrefixNotSpecified
				}
				project, err = query.Where("prefix = ?", prefix).First(db.Ctx)
				if err != nil {
					return err
				}
			} else if err != nil {
				return err
			}

			for i := range project.Issues {
				project.Issues[i].Project.GuildID = project.GuildID
			}

			// create a new state
			state := db.ProjectViewState{
				ProjectID: project.ID,
				Project:   project,
				Filter:    db.DefaultFilter(),
				Sorter:    db.DefaultSorter(),
			}

			return logic.InitIssueView(s, i, &state, false)
		}

		return nil
	},
}
