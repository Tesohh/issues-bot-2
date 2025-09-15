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
	environment := os.Getenv("DISCORD_ENVIRONMENT")

	if environment == "dev++" || environment == "dev" {
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
	} else {
		for _, c := range command.Commands {
			_, err := session.ApplicationCommandCreate(session.State.User.ID, "", &c.ApplicationCommand)
			if err != nil {
				return fmt.Errorf("Cannot create %s due to %s", c.Name, err.Error())
			}
		}
	}

	log.Println("Added commands")

	return nil
}
