package main

import (
	"fmt"
	"issues/v2/command"
	"issues/v2/slash"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var commands = map[string]*slash.Command{
	"project": &command.Project,
	"new":     &command.New,
}
var registeredCommands = make(map[string][]*discordgo.ApplicationCommand, 0)

func registerCommands(session *discordgo.Session) error {
	// TODO: add global commands in case of prod

	log.Println("Adding commands...")
	guildids := strings.SplitSeq(os.Getenv("DISCORD_GUILD_ID"), ",")

	for id := range guildids {
		registeredCommands[id] = make([]*discordgo.ApplicationCommand, 0)
		for _, c := range commands {
			cmd, err := session.ApplicationCommandCreate(session.State.User.ID, id, &c.ApplicationCommand)
			if err != nil {
				return fmt.Errorf("Cannot create %s due to %s", c.Name, err.Error())
			}
			registeredCommands[id] = append(registeredCommands[id], cmd)
		}
	}

	log.Println("Added commands")

	return nil
}
