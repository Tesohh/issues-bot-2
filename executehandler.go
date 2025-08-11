package main

import (
	"issues/v2/slash"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

var executeCommandHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if c, ok := commands[i.ApplicationCommandData().Name]; ok {
		if c.Disabled {
			slog.Warn("user somehow managed to execute a disabled command", slog.String("command", c.Name))
			return
		}
		err := c.Func(s, i.Interaction)
		if err != nil {
			slog.Error(err.Error())
			slash.ReplyWithEmbed(s, i.Interaction, discordgo.MessageEmbed{
				Title:       "Error",
				Description: err.Error(),
			}, true)
		}
	}
}
