package handler

import (
	"fmt"
	"issues/v2/command"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var RegisteredCommands = make(map[string][]*discordgo.ApplicationCommand, 0)

func RegisterCommands(session *discordgo.Session) error {
	// TODO: add global commands in case of prod

	log.Println("Adding commands...")
	guildids := strings.SplitSeq(os.Getenv("DISCORD_GUILD_ID"), ",")

	for id := range guildids {
		RegisteredCommands[id] = make([]*discordgo.ApplicationCommand, 0)
		for _, c := range command.Commands {
			cmd, err := session.ApplicationCommandCreate(session.State.User.ID, id, &c.ApplicationCommand)
			if err != nil {
				return fmt.Errorf("Cannot create %s due to %s", c.Name, err.Error())
			}
			RegisteredCommands[id] = append(RegisteredCommands[id], cmd)
		}
	}

	log.Println("Added commands")

	return nil
}
