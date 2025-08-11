package slash

import dg "github.com/bwmarrin/discordgo"

func GetOptionMapRaw(rawOptions []*dg.ApplicationCommandInteractionDataOption) map[string]*dg.ApplicationCommandInteractionDataOption {
	options := make(map[string]*dg.ApplicationCommandInteractionDataOption, len(rawOptions))
	for _, opt := range rawOptions {
		options[opt.Name] = opt
	}
	return options
}

func GetOptionMap(i *dg.Interaction) map[string]*dg.ApplicationCommandInteractionDataOption {
	rawOptions := i.ApplicationCommandData().Options
	return GetOptionMapRaw(rawOptions)
}
