package main

import (
	"fmt"
	"issues/v2/db"
	"issues/v2/handler"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
	"github.com/lmittmann/tint"
)

func initLogger() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))
}

func main() {
	initLogger()
	db.Connect(".data/issues2.db")

	session, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("DISCORD_BOT_TOKEN")))
	if err != nil {
		slog.Error(err.Error())
		return
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		slog.Info(fmt.Sprintf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator))
	})

	err = session.Open()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	slog.Info("adding handlers...")

	session.AddHandler(handler.GuildJoinHandler)
	session.AddHandler(handler.MessageDelete)
	session.AddHandler(handler.MessageCreate)
	session.AddHandler(handler.Router)

	slog.Info("registering commands...")
	err = handler.RegisterCommands(session)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if os.Getenv("DISCORD_ENVIRONMENT") == "dev" {
		log.Println("Removing commands...")
		for id, cmds := range handler.RegisteredCommands {
			for _, cmd := range cmds {
				err := session.ApplicationCommandDelete(session.State.User.ID, id, cmd.ID)
				if err != nil {
					slog.Error("Cannot delete '%v' command: %v", cmd.Name, err.Error())
				}
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
