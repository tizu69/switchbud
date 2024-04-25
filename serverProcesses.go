package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func killServer(channel string, v configServer) {
	_, err := s.ChannelMessageSendEmbed(channel,
		&discordgo.MessageEmbed{
			Title:       fmt.Sprintf(loc("process.kill"), v.DisplayName),
			Description: loc("process.kill.sub"),
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: v.DisplayIcon},
			Color:       hexToInt(loc("color.error")),
		})
	if err != nil {
		log.Warn("Failed to send kill message", "server", v.DisplayName, "err", err)
	}

	httpPost(ccEndpoint+"/servers/"+v.ServerID+"/action/stop_server", nil, ccToken)
	ccAwaitUserCount(v, -1, 600)
}

func bootServer(channel string, v configServer) {
	_, err := s.ChannelMessageSendEmbed(channel,
		&discordgo.MessageEmbed{
			Title:       fmt.Sprintf(loc("process.boot"), v.DisplayName),
			Description: loc("process.boot.sub"),
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: v.DisplayIcon},
			Color:       hexToInt(loc("color.generic")),
		})
	if err != nil {
		log.Warn("Failed to send boot message", "server", v.DisplayName, "err", err)
	}

	httpPost(ccEndpoint+"/servers/"+v.ServerID+"/action/start_server", nil, ccToken)

	if v.BootTime > 0 {
		time.AfterFunc(time.Duration(v.BootTime)*time.Second, func() {

			_, err := s.ChannelMessageSendEmbed(channel,
				&discordgo.MessageEmbed{
					Title:       fmt.Sprintf(loc("process.booted"), v.DisplayName),
					Description: loc("process.booted.sub"),
					Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: v.DisplayIcon},
					Color:       hexToInt(loc("color.success")),
				})
			if err != nil {
				log.Warn("Failed to send boot timer message", "server", v.DisplayName, "err", err)
			}
		})
	}
}
