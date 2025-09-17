package command

import (
	"issues/v2/db"
	"issues/v2/slash"
	"os"

	dg "github.com/bwmarrin/discordgo"
)

var Reset = slash.Command{
	ApplicationCommand: dg.ApplicationCommand{
		Name:        "reset",
		Description: "Deletes everything about this guild: projects, issues, roles. after this, kick the bot and readd it.",
		Options: []*dg.ApplicationCommandOption{
			{
				Type:        dg.ApplicationCommandOptionBoolean,
				Name:        "confirm",
				Description: "are you really sure?",
				Required:    true,
			},
			{
				Type:        dg.ApplicationCommandOptionString,
				Name:        "password",
				Description: "password only known by tesohh",
				Required:    true,
			},
		},
	},
	Func: func(s *dg.Session, i *dg.Interaction) error {

		options := slash.GetOptionMap(i)
		confirm := options["confirm"].BoolValue()
		password := options["password"].StringValue()
		if !confirm {
			return slash.ReplyWithEmbed(s, i, dg.MessageEmbed{
				Title: "alright, no actions taken",
			}, true)
		}

		if password != os.Getenv("DISCORD_RESET_PASSWORD") {
			return slash.ReplyWithEmbed(s, i, dg.MessageEmbed{
				Title: "wrong password. no actions taken",
			}, true)
		}

		guild, err := db.Guilds.
			Preload("Roles", nil).
			Preload("Projects", nil).
			Where("id = ?", i.GuildID).
			First(db.Ctx)
		if err != nil {
			return err
		}

		_, err = db.Guilds.Where("id = ?", i.GuildID).Delete(db.Ctx)
		if err != nil {
			return err
		}

		for _, project := range guild.Projects {
			err = ProjectDelete(s, i, project.Prefix, confirm, false)
			if err != nil {
				return err
			}
		}

		for _, role := range guild.Roles {
			err = s.GuildRoleDelete(i.GuildID, role.ID)
			if err != nil {
				return err
			}
		}

		_, err = db.Roles.Where("guild_id = ?", i.GuildID).Delete(db.Ctx)
		if err != nil {
			return err
		}

		return slash.ReplyWithEmbed(s, i, dg.MessageEmbed{Title: "Harris..."}, false)
	},
}
