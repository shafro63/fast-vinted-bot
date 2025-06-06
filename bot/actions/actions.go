package actions

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"fast-vinted-bot/apicalls"
	"fast-vinted-bot/bot/components"
	"fast-vinted-bot/cache"
	"fast-vinted-bot/database"
	"fast-vinted-bot/services"
	"fast-vinted-bot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load()

// Fetch Items and Send them on your discord channel
func FetchAndSendToDiscord(s *discordgo.Session, data *utils.DiscordUserData, rb *utils.RequestBuilder) {
	duration, err := strconv.Atoi(os.Getenv("FETCHING_DURATION"))
	if err != nil {
		slog.Error("unable to get duration entry", "error", err)
		os.Exit(1)
	}

	timer := &cache.Timer{
		TickerChannel: make(chan struct{}),
		Duration:      time.Duration(duration) * time.Hour,
	}
	cache.TimerCache.SetTimer(data.LinkName, timer)

	dataChan := make(chan []utils.CatalogItem)
	stopChan := make(chan bool)
	var lastId int64

	cache.DataCache.SetMonitorSession(data, stopChan)

	apicalls.RefreshCookie(rb)
	services.FetchCatalogAtInterval(rb, timer, dataChan, stopChan)

	// Get the latest items
	for items := range dataChan {
		latest_items := services.LatestItems(items, &lastId)
		if latest_items == nil {
			continue
		}

		for i := range len(latest_items) {
			item := latest_items[i]

			embedItem := components.CreateEmbed(&item, data.LinkName)
			button := components.CreateActionsRow(&item)

			msg := &discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{
					embedItem,
				},
				Components: []discordgo.MessageComponent{
					button,
				},
			}
			_, err = s.ChannelMessageSendComplex(data.ChannelID, msg)
			if err != nil {
				slog.Error("Error while sending ads to discord", "channel", data.ChannelID, "error", err)
			}

		}

	}
}

// This is launched at startup and load the sessions
func LaunchMonitorSessions(s *discordgo.Session) {
	var wg sync.WaitGroup

	channels, err := database.GetAllActiveChannels("db", "vinted")
	if err != nil {
		slog.Error("unable to launch monitor sessions", "error", err)
		os.Exit(1)
	}

	for _, v := range channels {
		userChannels := v.Channels
		wg.Add(1)

		go func(userchannels map[string]*utils.ChannelInfo) {
			defer wg.Done()

			for k, v := range userChannels {
				wg.Add(1)

				go func(links map[string]string, channelID string) {
					defer wg.Done()

					for linkName, link := range links {
						go func(linkname string, link string) {
							data := &utils.DiscordUserData{
								ChannelID: channelID,
								LinkName:  linkname,
								Link:      link,
							}

							parsedUrl, err := services.ParsedUrl(data.Link)
							if err != nil {
								msgError := fmt.Sprintf("Your link %v is invalid", link)
								slog.Error("corrupted link", "channel", channelID, "link", link, "error", err)
								s.ChannelMessageSend(channelID, msgError)
								return
							}

							rb := utils.NewRequestBuilder()
							rb.Method = "GET"
							rb.URL = parsedUrl

							c := apicalls.GetCookie(rb)
							cookie := apicalls.FormatedAuthCookie(c)
							if cookie == nil {
								slog.Error("nil authcookie")
								return
							}
							rb.Cookie = cookie

							FetchAndSendToDiscord(s, data, rb)

						}(linkName, link)
					}
				}(v.Links, k)
			}
		}(userChannels)
	}
	wg.Wait()
	slog.Info("all monitor sessions set up !")
}
