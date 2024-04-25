package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func dMainCommand(s *discordgo.Session, i *discordgo.InteractionCreate, cfg config) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	embeds := []*discordgo.MessageEmbed{}
	runningServers := ccGetRunningServers(cfg.Servers)
	for _, v := range runningServers {
		embeds = append(embeds, getServerEmbed(v))
	}

	embeds = append(embeds, &discordgo.MessageEmbed{
		Title: func() string {
			if len(embeds) > 0 {
				return loc("selector.titleSwitch")
			}
			return loc("selector.title")
		}(),

		Description: loc("selector.desc"),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf(loc("selector.count"), len(cfg.Servers)),
		},
		Color: hexToInt(loc("color.generic")),
	})

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
		Components: bSingleComp(discordgo.SelectMenu{
			MenuType:    discordgo.StringSelectMenu,
			CustomID:    "serverSelect",
			Placeholder: loc("selector.selectPlaceholder"),

			Options: func() []discordgo.SelectMenuOption {
				var options []discordgo.SelectMenuOption
				for _, v := range cfg.Servers {
					options = append(options, discordgo.SelectMenuOption{
						Label: v.DisplayName,
						Value: v.ServerID,
						Emoji: bEmojiFromUrl(v.DisplayIcon, cfg.GuildId),
					})
				}
				return options
			}(),
		}),
	})
}

func getServerEmbed(v configServer) *discordgo.MessageEmbed {
	onlineCount := ccGetPlayerCount(v)

	var color int
	if onlineCount >= 0 {
		color = hexToInt(loc("color.success"))
	} else {
		color = hexToInt(loc("color.generic"))
	}

	return &discordgo.MessageEmbed{
		Title:       v.DisplayName,
		Description: v.DisplayDesc,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: v.DisplayIcon},
		Footer: &discordgo.MessageEmbedFooter{
			Text: func() string {
				if onlineCount < 0 {
					return loc("pack.onlineCount.empty")
				}
				return fmt.Sprintf(loc("pack.onlineCount"), onlineCount)
			}(),
		},
		Color: color,
	}
}
