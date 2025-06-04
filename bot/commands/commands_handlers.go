package commands

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"fast-vinted-bot/cache"
	"fast-vinted-bot/database"
	"fast-vinted-bot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	_ = godotenv.Load()

	ChannelCommandID = os.Getenv("DISCORD_COMMAND_CHANNEL_ID")

	// Handlers for discord commands
	// See the discordgo api for additionnal info
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){

		"create_private_channel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			data := &utils.DiscordUserData{
				GuildID:     i.GuildID,
				Member:      i.Member,
				ChannelName: i.ApplicationCommandData().Options[0].StringValue(),
			}
			user := data.Member.User

			// Fast reply (mandatory)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⏳ Creating private channel...",
					Flags:   discordgo.MessageFlagsEphemeral, // ephemeral (visible seulement pour l'utilisateur)
				},
			})

			err := database.CreateUser("db", "vinted", data)
			if err != nil {
				msgError := "❌ Can't create user data"
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &msgError,
				})
				slog.Error("Create user error", "user", user.ID, "error", err)
				return
			}

			channels, err := database.GetChannels("db", "vinted", data)
			if err != nil {
				slog.Error("Get Channels error", "user", user.ID, "error", err)
				msgError := "❌ Can't get user data"
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &msgError,
				})
				return
			}

			max_channels, err := strconv.Atoi(os.Getenv("MAX_PRIVATE_CHANNELS"))
			if err != nil {
				slog.Error("unable to get env entry", "error", err)
				os.Exit(1)
			}

			if len(channels) >= max_channels {
				msgError := "❌ You have reached your private channels limit"
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &msgError,
				})
				return
			}

			// Create the private channel
			privateChannel, err := s.ThreadStartComplex(ChannelCommandID, &discordgo.ThreadStart{
				Name: data.ChannelName,
				Type: discordgo.ChannelTypeGuildText,
			})
			slog.Debug("Create Private channel Success", "user", user.ID, "name", data.ChannelName, "id", privateChannel.ID)
			if err != nil {
				slog.Error("Create Private channel error", "user", user.ID, "name", data.ChannelName, "id", privateChannel.ID)

				msgError := "❌ Unable to create private channel, please inform the staff if it persists."
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &msgError,
				})
				return
			}

			data.ChannelID = privateChannel.ID

			err = database.SetChannel("db", "vinted", data)
			if err != nil {
				slog.Error("Channel Set error", "user", user.ID, "error", err)

				_, ErrDelete := s.ChannelDelete(data.ChannelID)

				slog.Error("Channel delete error", "user", user.ID, "error", ErrDelete)

				msgError := "❌ Unable to create private channel, please inform the admins if the problem persists."
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &msgError,
				})

				return
			}
			slog.Debug("Channel Set", "user", user.ID, "channelID", data.ChannelID)

			err = s.ThreadMemberAdd(data.ChannelID, user.ID)
			if err != nil {
				slog.Error("Thread member Add Error", "user", user.ID, "channelID", data.ChannelID)

				msgError := "❌ Unable to create private channel, please inform admins if it persists."
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &msgError,
				})
				return
			}

			// 💬 Envoie un message dans le channel créé
			s.ChannelMessageSend(data.ChannelID, fmt.Sprintf("🎫 @%s Welcome to your private channel :D", user.Username))

		},

		"add_link": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⏳Processing...",
					Flags:   discordgo.MessageFlagsEphemeral, // ephemeral (visible seulement pour l'utilisateur)
				},
			})

			values := i.ApplicationCommandData().Options
			data := &utils.DiscordUserData{
				Member:   i.Member,
				LinkName: values[0].StringValue(),
				Link:     values[1].StringValue(),
			}
			user := data.Member.User

			user_channels, err := database.GetChannels("db", "vinted", data)
			if err != nil {
				slog.Debug("Get Channels Error", "error", err)
				msgError := fmt.Sprintf("%v", err)
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &msgError,
				})
				return
			} else if len(user_channels) == 0 {
				msgError := "You should create a private channel first"
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &msgError,
				})
				return
			}

			cache.DataCache.SetUserData(user.ID, data)
			slog.Debug("Cache data set", "user", user.ID, "channelID", data.ChannelID)

			var channels_menu []discordgo.SelectMenuOption
			for k, v := range user_channels {
				choices := discordgo.SelectMenuOption{
					Label:       v.Name,
					Value:       k,
					Description: "",
				}
				channels_menu = append(channels_menu, choices)
			}

			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Components: &[]discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								MenuType:    discordgo.StringSelectMenu,
								Placeholder: "Select a channel",
								CustomID:    "add_link_menu",
								Options:     channels_menu,
								MaxValues:   1,
							},
						},
					},
				},
			})
			if err != nil {
				slog.Error("Interaction response edit error", "user", user.ID, "error", err, "channelID", data.ChannelID)
			}
		},

		"delete_private_channel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⏳Processing...",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

			data := &utils.DiscordUserData{
				Member: i.Member,
			}
			user := data.Member.User

			channels, err := database.GetChannels("db", "vinted", data)
			if err != nil {
				slog.Debug("Get channels error", "user", user.ID, "error", err)
				return
			}

			var channels_menu []discordgo.SelectMenuOption
			for k, v := range channels {
				choices := discordgo.SelectMenuOption{
					Label:       v.Name,
					Value:       k,
					Description: "You can select only one",
				}
				channels_menu = append(channels_menu, choices)
			}

			msg := "Select a channel"
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &msg,
				Components: &[]discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								MenuType:    discordgo.StringSelectMenu,
								Placeholder: "Select a channel to delete",
								CustomID:    "delete_channel_menu",
								Options:     channels_menu,
								MaxValues:   1,
							},
						},
					},
				},
			})
			if err != nil {
				slog.Error("Interaction respond error", "user", user.ID, "error", err)
				return
			}

		},

		"delete_link": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⏳Processing...",
				},
			})

			data := &utils.DiscordUserData{
				Member: i.Member,
			}
			user := data.Member.User

			userChannels, err := database.GetChannels("db", "vinted", data)
			if err != nil {
				msg := fmt.Sprintf("%v", err)
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &msg,
				})
				slog.Debug("Get channels error", "user", user.ID, "error", err)
				return
			} else if len(userChannels) == 0 {
				msg := "❌ You must create a private channel first"
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &msg,
				})
				return
			}

			var channels_menu []discordgo.SelectMenuOption
			for k, v := range userChannels {
				choices := discordgo.SelectMenuOption{
					Label:       v.Name,
					Value:       k,
					Description: "You can select only one",
				}
				channels_menu = append(channels_menu, choices)
			}

			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Components: &[]discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								MenuType:    discordgo.StringSelectMenu,
								Placeholder: "Select a channel to delete link from",
								CustomID:    "delete_link_channel_menu",
								Options:     channels_menu,
								MaxValues:   1,
							},
						},
					},
				},
			})
			if err != nil {
				slog.Error("Interaction response edit error", "user", user.ID, "error", err, "channelID", data.ChannelID)
			}
		},
	}
)
