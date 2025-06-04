package commands

import (
	"github.com/bwmarrin/discordgo"
)

// Discord commands names
var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "create_private_channel",
			Description: "Create a private channel",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "Private channel name",
				Required:    true,
			},
			},
		},
		{
			Name:        "delete_private_channel",
			Description: "Delete a private channel",
			Type:        discordgo.ChatApplicationCommand,
		},
		{
			Name:        "add_link",
			Description: "Add a vinted link to a private channel",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "Unique link name",
				Required:    true,
			},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "link",
					Description: "Link",
					Required:    true,
				},
			},
		},
		{
			Name:        "delete_link",
			Description: "Delete a link from a private channel",
			Type:        discordgo.ChatApplicationCommand,
		},
	}
)
