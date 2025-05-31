package components

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"fast-vinted-bot/utils"

	"github.com/bwmarrin/discordgo"
)

func hexToInt(hex string) int {
	hex = strings.TrimPrefix(hex, "#")          // enlève le '#' s'il y en a un
	value, err := strconv.ParseInt(hex, 16, 32) // base 16 = hexadécimal
	if err != nil {
		return 0
	}
	return int(value)
}

func CreateEmbed(item utils.CatalogItem, name string) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:     item.User.ProfileURL,
			Name:    item.User.Login,
			IconURL: item.User.Photo.URL,
		},
		URL:   item.URL,
		Title: item.Title,
		Color: hexToInt(item.Photo.DominantColor),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "🕒Updated",
				Value:  fmt.Sprintf("<t:%d:R>", time.Now().Unix()),
				Inline: true,
			},
			{
				Name:   "💳Price",
				Value:  item.Price.Amount + "" + item.Price.CurrencyCode,
				Inline: true,
			},
			{
				Name:   "📏Size",
				Value:  item.SizeTitle,
				Inline: true,
			},
			{
				Name:   "🏷️Brand",
				Value:  item.BrandTitle,
				Inline: true,
			},
			{
				Name:   "🎭Status",
				Value:  item.Status,
				Inline: true,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: item.Photo.URL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: name,
		},
	}
	return embed
}
