package command

import (
	"issues/v2/dataview"
	"issues/v2/db"
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

			issues := project.Issues
			for i := range issues {
				issues[i].Project.GuildID = project.GuildID
			}

			// create a new state
			state := db.ProjectViewState{
				ProjectID: project.ID,
				Project:   project,
				Filter:    db.DefaultFilter(),
				Sorter:    db.DefaultSorter(),
			}
			// send the message with no buttons
			components := dataview.MakeInteractiveIssuesView(issues, &state, true)

			err = slash.ReplyWithComponents(s, i, false, components...)
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
			err = db.ProjectViewStates.Create(db.Ctx, &state)
			if err != nil {
				return err
			}

			state.Project = project // im so sorry
			// make the view WITH the buttons
			components = dataview.MakeInteractiveIssuesView(issues, &state, false)
			_, err = s.InteractionResponseEdit(i, &dg.WebhookEdit{
				Components: &components,
			})
			if err != nil {
				return err
			}

			// TODO: apply filters from the discord options
		}

		return nil
	},
}
