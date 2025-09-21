package command

import (
	"fmt"
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
					{
						Type:        dg.ApplicationCommandOptionBoolean,
						Name:        "detached",
						Description: "if true, the message will be sent as a message and not a reply",
					},
				},
			},
			{
				Type:        dg.ApplicationCommandOptionSubCommand,
				Name:        "projects",
				Description: "lists all projects in this guild",
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
				Preload("Issues", func(query gorm.PreloadBuilder) error {
					query.Where("kind = ?", db.IssueKindNormal)
					return nil
				}).
				Preload("Issues.Tags", nil).
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

			detached := false
			if detachedOpt, ok := options["detached"]; ok {
				detached = detachedOpt.BoolValue()
			}
			if detached {
				slash.ReplyWithText(s, i, "loading", true)
				s.InteractionResponseDelete(i)

				err = logic.InitIssueViewDetached(s, i, i.ChannelID, &state)
				if err != nil {
					return err
				}
			} else {
				err = logic.InitIssueView(s, i, &state, false)
				if err != nil {
					return err
				}
			}

			// solicitation
			count, err := db.ProjectViewStates.Where("project_id = ?", project.ID).Count(db.Ctx, "*")
			if count >= 5 && !project.HasBeenSolicitedByListWarning {
				s.ChannelMessageSendComplex(i.ChannelID, &dg.MessageSend{
					Components: []dg.MessageComponent{
						slash.StandardizeContainer(warningContainer),
					},
					Flags: dg.MessageFlagsIsComponentsV2,
				})
				db.Projects.Where("id = ?", project.ID).Update(db.Ctx, "HasBeenSolicitedByListWarning", true)
			}

			return nil
		case "projects":
			projects, err := db.Projects.Where("guild_id = ?", i.GuildID).Find(db.Ctx)
			if err != nil {
				return err
			}

			str := ""
			for _, project := range projects {
				str += fmt.Sprintf("- `[%s]` %s <#%s>\n", strings.ToUpper(project.Prefix), project.Name, project.IssuesInputChannelID)
			}

			embed := dg.MessageEmbed{
				Title:       "Projects in this guild",
				Description: str,
			}
			return slash.ReplyWithEmbed(s, i, embed, false)
		}

		return nil
	},
}

var warningContainer = dg.Container{
	Components: []dg.MessageComponent{
		dg.TextDisplay{Content: `
# ⚠️ Warning 
- Having too many lists can get really slow due to discord rate limits.
  - This means that lists are updated sequentially when you add a new issue, update it etc.
  - The first 5 or so lists are updated almost instantly, while it takes about 5 seconds or more for each next list.
  - The AutoList is always the first one to be updated.
- After you're done using a temporary list, please delete the message.
- Old lists are automatically untracked after a day, unless you set your lists to` + " `Permanent`" + `.
- This won't be shown again in this project, but be wary.`},
	},
}
