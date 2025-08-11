package slash

import "github.com/bwmarrin/discordgo"

type Command struct {
	discordgo.ApplicationCommand
	Disabled bool
	Func     func(s *discordgo.Session, i *discordgo.Interaction) error
}
