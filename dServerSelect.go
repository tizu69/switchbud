package main

import (
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func dServerSelect(s *discordgo.Session, i *discordgo.InteractionCreate, cfg config) {
	id := i.MessageComponentData().Values[0]
	server := slices.IndexFunc(cfg.Servers, func(v configServer) bool {
		return v.ServerID == id
	})
	if server == -1 {
		log.Warn("Failed to find server", "server", id)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: loc("generic.failed"),
			},
		})
		return
	}
	embed := getServerEmbed(cfg.Servers[server])

	errError(s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{
					Components: func() []discordgo.MessageComponent {
						var buttons []discordgo.MessageComponent

						if ccGetPlayerCount(cfg.Servers[server]) > 0 {
							buttons = append(buttons, bButton(loc("pack.switch.occupied"),
								discordgo.SecondaryButton, "?", true))
						} else if ccGetPlayerCount(cfg.Servers[server]) == 0 {
							buttons = append(buttons, bButton(loc("pack.switch.stop"),
								discordgo.DangerButton, "switchTo/"+cfg.Servers[server].ServerID, false))
						} else {
							buttons = append(buttons, bButton(loc("pack.switch"),
								discordgo.PrimaryButton, "switchTo/"+cfg.Servers[server].ServerID, false))
						}

						if cfg.Servers[server].DownloadLink != "" {
							buttons = append(buttons, discordgo.Button{
								Label: loc("pack.download"),
								Style: discordgo.LinkButton,
								URL:   cfg.Servers[server].DownloadLink,
							})
						}

						return buttons
					}(),
				},
			},
		},
	}))
}
