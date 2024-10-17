package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"

	"github.com/chand1012/hooky/pkg/config"
	"github.com/chand1012/hooky/pkg/handle"
)

var (
	Token     = flag.String("token", "", "Bot authentication token")
	App       = flag.String("app", "", "Application ID")
	Guild     = flag.String("guild", "", "Guild ID")
	ConfigDir = flag.String("config-dir", "./config", "Configuration directory")
)

func main() {
	flag.Parse()
	if *App == "" {
		app_id := os.Getenv("APP_ID")
		if app_id != "" {
			App = &app_id
		} else {
			log.Fatal("application id is not set. Please provide it with -app or environment variable APP_ID")
		}
	}

	// do the same for token. Guild is optional
	if *Token == "" {
		token := os.Getenv("BOT_TOKEN")
		if token != "" {
			Token = &token
		} else {
			log.Fatal("bot token is not set. Please provide it with -token or environment variable BOT_TOKEN")
		}
	}

	if *Guild == "" {
		guild := os.Getenv("GUILD_ID")
		if guild != "" {
			Guild = &guild
		}
	}

	cmds, err := config.LoadConfigs(*ConfigDir)
	if err != nil {
		log.Fatalf("could not load configs: %s", err)
	}

	session, _ := discordgo.New("Bot " + *Token)

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		handle.Command(cmds, s, i)
	})

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Infof("Logged in as %s", r.User.String())
	})

	commands, err := cmds.ToApplicationCommands()
	if err != nil {
		log.Fatalf("could not convert commands: %s", err)
	}

	_, err = session.ApplicationCommandBulkOverwrite(*App, *Guild, commands)
	if err != nil {
		log.Fatalf("could not register commands: %s", err)
	}

	err = session.Open()
	if err != nil {
		log.Fatalf("could not open session: %s", err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = session.Close()
	if err != nil {
		log.Errorf("could not close session gracefully: %s", err)
	}
}
