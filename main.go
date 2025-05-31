package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"

	"fast-vinted-bot/bot/actions"
	"fast-vinted-bot/bot/commands"
	"fast-vinted-bot/cache"
	"fast-vinted-bot/database"
	"fast-vinted-bot/logger"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	_              = godotenv.Load()
	GuildID        = flag.String("guild", os.Getenv("DISCORD_GUILD_ID"), "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", os.Getenv("DISCORD_BOT_TOKEN"), "token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	ClientID       = flag.String("client", os.Getenv("DISCORD_CLIENT_ID"), "ClientID")
	s              *discordgo.Session
)

func init() { flag.Parse() }

func init() { logger.InitLogger() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		slog.Error("Invalid bot parameters", "error", err)
		os.Exit(1)
	}
}

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}

		case discordgo.InteractionMessageComponent:
			if h, ok := commands.MsgComponentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	slog.Info("Handlers initialised")
}

func main() {

	database.InitMongoClient()
	defer func() {
		if err := database.Client.Disconnect(context.Background()); err != nil {
			slog.Error("Unable to disconnect from MongoDB", "error", err)
		} else {
			slog.Info("🔌 Successfully disconnected from MongoDB.")
		}
	}()

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		slog.Info("Logged in", "username", s.State.User.Username+s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		slog.Error("Cannot open the session", "error", err)
		os.Exit(1)
	}

	slog.Info("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands.Commands))
	for i, v := range commands.Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			slog.Error("Cannot create command", "command", v.Name, "error", err)
			os.Exit(1)
		}
		registeredCommands[i] = cmd
	}

	cache.LaunchTicker()
	actions.LaunchMonitorSessions(s)

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	slog.Info("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		slog.Info("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				slog.Error("Cannot delete command", "command", v.Name)
				os.Exit(1)

			}
		}
	}

	slog.Info("Gracefully shutting down.")

}
