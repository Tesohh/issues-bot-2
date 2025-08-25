package slash

import dg "github.com/bwmarrin/discordgo"

func standardizeContainer(container dg.Container) dg.Container {
	if container.AccentColor == nil || *container.AccentColor == 0 {
		container.AccentColor = Ptr(EmbedColor)
	}

	return container
}

func ReplyWithComponents(s *dg.Session, i *dg.Interaction, ephemeral bool, components ...dg.MessageComponent) error {
	var flags dg.MessageFlags = dg.MessageFlagsIsComponentsV2
	if ephemeral {
		flags = dg.MessageFlagsEphemeral
	}

	for i := range components {
		if container, ok := components[i].(dg.Container); ok {
			components[i] = standardizeContainer(container)
		}
	}

	return s.InteractionRespond(i, &dg.InteractionResponse{
		Type: dg.InteractionResponseChannelMessageWithSource,
		Data: &dg.InteractionResponseData{
			Components: components,
			Flags:      flags,
		},
	})
}
