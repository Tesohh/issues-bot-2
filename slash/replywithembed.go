package slash

import (
	"github.com/bwmarrin/discordgo"
)

func ReplyWithText(s *discordgo.Session, i *discordgo.Interaction, content string, ephemeral bool) error {
	var flags discordgo.MessageFlags
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}

	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   flags,
		},
	})
}

func ReplyWithEmbed(s *discordgo.Session, i *discordgo.Interaction, embed discordgo.MessageEmbed, ephemeral bool) error {
	embed = standardizeEmbed(embed)

	var flags discordgo.MessageFlags
	if ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}

	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
			Flags:  flags,
		},
	})
}
