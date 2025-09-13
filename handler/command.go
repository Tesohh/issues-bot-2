package handler

import (
	"issues/v2/command"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func Command(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if c, ok := command.Commands[i.ApplicationCommandData().Name]; ok {
		if c.Disabled {
			slog.Warn("user somehow managed to execute a disabled command", slog.String("command", c.Name))
			return
		}
		err := c.Func(s, i.Interaction)
		if err != nil {
			slog.Error(err.Error())
			embed := discordgo.MessageEmbed{
				Title:       "Error",
				Description: err.Error(),
				Color:       0xFF0000,
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{&embed},
					Flags:  discordgo.MessageFlagsEphemeral,
				},
			})
		}
	}
}
