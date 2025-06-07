package commands

import (
	"fast-vinted-bot/apicalls"
	"fast-vinted-bot/bot/actions"
	"fast-vinted-bot/cache"
	"fast-vinted-bot/database"
	"fast-vinted-bot/services"
	"fast-vinted-bot/utils"
	"log/slog"

	"fmt"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	_ = godotenv.Load()

	// Handlers for discord message components
	MsgComponentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"add_link_menu": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.Member.User
			values := i.MessageComponentData().Values
			data := cache.DataCache.GetUserData(user.ID)
			data.ChannelID = values[0]

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
				Data: &discordgo.InteractionResponseData{
					Content: "⏳Adding link...",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

			userlinks, err := database.GetLinks("db", "vinted", data)
			if err != nil {
				slog.Debug("Get links error", "user", user.ID, "error", err, "channelID", data.ChannelID)

				msgError := fmt.Sprintf("%v", err)
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content:    &msgError,
					Components: &[]discordgo.MessageComponent{},
				})
				return
			}

			max_links, err := strconv.Atoi(os.Getenv("MAX_LINKS_PER_CHANNEL"))
			if err != nil {
				slog.Error("unable to get env entry", "error", err)
				os.Exit(1)
			}

			if len(userlinks) >= max_links {
				msgError := "You have reached your links limit in this channel"
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content:    &msgError,
					Components: &[]discordgo.MessageComponent{},
				})
				return
			} else if _, ok := userlinks[data.LinkName]; ok {
				msgError := "This link's name already exists. Please retry with a different name"
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content:    &msgError,
					Components: &[]discordgo.MessageComponent{},
				})
				return
			}

			parsedUrl, err := services.ParsedUrl(data.Link)
			if err != nil {
				s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("%v", err))
				return
			}

			err = database.SetLink("db", "vinted", data)
			if err != nil {
				slog.Error("Set link error", "user", user.ID, "error", err, "channelID", data.ChannelID)

				msgError := fmt.Sprintf("%v", err)
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content:    &msgError,
					Components: &[]discordgo.MessageComponent{},
				})
				return
			}

			msgSuccess := fmt.Sprintf("✅ Link added to your channel <#%v>", data.ChannelID)
			s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content:    &msgSuccess,
				Components: &[]discordgo.MessageComponent{},
			})

			rb := utils.Rb
			rb.Method = "GET"
			rb.URL = parsedUrl

			c := apicalls.GetCookie(rb)
			cookie := apicalls.FormatedAuthCookie(c)
			if cookie == nil {
				s.ChannelMessageSend(i.ChannelID, "Can't get vinted cookie, please contact admins")
				slog.Debug("Can't get vinted cookie")
				return
			}
			rb.Cookie = cookie

			actions.FetchAndSendToDiscord(s, data, rb)
		},
		"delete_link_channel_menu": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			values := i.MessageComponentData().Values
			data := &utils.DiscordUserData{
				Member:    i.Member,
				ChannelID: values[0],
			}
			user := data.Member.User

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
				Data: &discordgo.InteractionResponseData{
					Content: "⏳Processing...",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

			userLinks, err := database.GetLinks("db", "vinted", data)
			if err != nil {
				msgError := fmt.Sprintf("%v", err)
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content:    &msgError,
					Components: &[]discordgo.MessageComponent{},
				})
				return
			} else if len(userLinks) == 0 {
				msgError := "❌This channel doesn't contain any link yet"
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content:    &msgError,
					Components: &[]discordgo.MessageComponent{},
				})
				return
			}

			cache.DataCache.SetUserData(user.ID, data)
			slog.Debug("Data Cache Links Set")

			var links_menu []discordgo.SelectMenuOption
			for k, v := range userLinks {
				choices := discordgo.SelectMenuOption{
					Label:       k,
					Value:       k,
					Description: v[:100],
				}
				links_menu = append(links_menu, choices)
			}

			msg := "select"
			_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &msg,
				Components: &[]discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.SelectMenu{
								MenuType:    discordgo.StringSelectMenu,
								Placeholder: "Select a link to delete",
								CustomID:    "delete_link_menu",
								Options:     links_menu,
								MaxValues:   1,
							},
						},
					},
				},
			})
			if err != nil {
				slog.Error("message not sent", "error", err)
			}
		},

		"delete_link_menu": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
				Data: &discordgo.InteractionResponseData{
					Content: "⏳Processing...",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

			values := i.MessageComponentData().Values
			data := &utils.DiscordUserData{
				Member:   i.Member,
				LinkName: values[0],
			}
			user := data.Member.User

			cachedata := cache.DataCache.GetUserData(user.ID)
			if cachedata == nil {
				return
			}
			slog.Debug("Cache Get data success")

			data.ChannelID = cachedata.ChannelID

			err := database.DeleteLink("db", "vinted", data)
			if err != nil {
				slog.Error("Delete Link Error ", "user", user.ID, "error", err)

				msgError := fmt.Sprintf("%v", err)
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &msgError,
				})
				return
			}

			stopChan := cache.DataCache.GetMonitorSession(data)
			if stopChan == nil {
				msgSuccess := "✅Link properly deleted from your channel"
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Content:    msgSuccess,
						Flags:      discordgo.MessageFlagsEphemeral,
						Components: []discordgo.MessageComponent{},
					},
				})
				return
			} else {

				cache.TimerCache.DeleteTimer(data.LinkName)
				cache.DataCache.DeleteMonitorSession(data)

				msgSuccess := "✅ Link properly deleted from your channel"
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Content:    msgSuccess,
						Flags:      discordgo.MessageFlagsEphemeral,
						Components: []discordgo.MessageComponent{},
					},
				})
			}

		},
		"delete_channel_menu": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
				Data: &discordgo.InteractionResponseData{
					Content: "⏳Processing...",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})

			values := i.MessageComponentData().Values
			data := &utils.DiscordUserData{
				Member:    i.Member,
				ChannelID: values[0],
			}

			err := database.DeleteChannel("db", "vinted", data)
			if err != nil {
				slog.Error("Delete Channel Error", "database", "db", "user", data.Member.User.ID, "channel", data.ChannelID, "error", err)
				msgError := fmt.Sprintf("%v", err)
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content:    &msgError,
					Components: &[]discordgo.MessageComponent{},
				})
				return
			}

			cache.TimerCache.DeleteAllTimersInChannel(data)
			cache.DataCache.DeleteMonitoringChannel(data)

			_, err = s.ChannelDelete(data.ChannelID)
			if err != nil {
				msgError := fmt.Sprintf("%v", err)
				s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content:    &msgError,
					Components: &[]discordgo.MessageComponent{},
				})
				return
			}

			_, err = s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Content: "✅ Private channel properly deleted",
				Flags:   discordgo.MessageFlagsEphemeral,
			})
			if err != nil {
				slog.Error("can't send message", "database", "db", "user", data.Member.User.ID, "error", err)
				return
			}

		},
	}
)
