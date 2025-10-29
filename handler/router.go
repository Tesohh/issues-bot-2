package handler

import (
	"github.com/bwmarrin/discordgo"
)

func Router(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		Command(s, i)
	case discordgo.InteractionMessageComponent:
		MessageComponent(s, i)
	case discordgo.InteractionModalSubmit:
		ModalSubmit(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		Autocomplete(s, i)
	}
}
