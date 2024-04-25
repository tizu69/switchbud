package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func dSwitchTo(s *discordgo.Session, i *discordgo.InteractionCreate, cfg config) {
	id := strings.Split(i.MessageComponentData().CustomID, "/")[1]
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
	srv := cfg.Servers[server]

	currentPlayerCount := ccGetPlayerCount(srv)
	if currentPlayerCount == 0 {
		s.InteractionResponseDelete(i.Interaction)
		killServer(i.ChannelID, srv)

		_, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: fmt.Sprintf(loc("pack.shutDown"), srv.DisplayName),
		})
		if err != nil {
			log.Warn("Failed to send shutdown message", "server", srv.DisplayName, "err", err)
		}

		return
	} else if currentPlayerCount > 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Description: loc("pack.switch.occupied"),
						Color:       hexToInt(loc("color.error")),
					},
				},
			},
		})
		return
	}

	runningEmpty, resourceSlotsUnused := ccGetRunningEmptyServers(cfg.Servers)
	toKill := serverPurge(runningEmpty, srv.ResourceSlots, cfg.ResourceSlots-resourceSlotsUnused)

	if toKill == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Description: loc("pack.switchingFailed"),
						Color:       hexToInt(loc("color.error")),
					},
				},
			},
		})
		return
	}

	errError(s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Description: loc("pack.switching"),
					Color:       hexToInt(loc("color.success")),
				},
			},
		},
	}))

	for _, v := range toKill {
		killServer(i.ChannelID, v)
	}

	bootServer(i.ChannelID, srv)
}
