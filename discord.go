package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

var s *discordgo.Session

func setupDiscord(cfg config, removeCommands bool) {
	dg, err := discordgo.New(cfg.DiscordToken)
	if err != nil {
		log.Fatal("Failed to create Discord session", "err", err)
	}
	s = dg

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Info("Authenticated with Discord", "user", s.State.User.Username)

		if removeCommands {
			_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID,
				cfg.GuildId, []*discordgo.ApplicationCommand{})
			if err != nil {
				log.Error("Failed to remove commands", "err", err)
			}
		}

		_, err := s.ApplicationCommandCreate(s.State.User.ID, cfg.GuildId, &discordgo.ApplicationCommand{
			Name:         loc("cmd.packs.title"),
			Description:  loc("cmd.packs.description"),
			DMPermission: pointer(false),
		})
		if err != nil {
			log.Error("Failed to create slash command", "err", err)
		}
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i, cfg)
				return
			}
			log.Warn("Failed to handle interaction", "ref", i.ApplicationCommandData().Name)
		case discordgo.InteractionMessageComponent:
			if h, ok := commandHandlers[strings.Split(
				i.MessageComponentData().CustomID, "/")[0]+"[comp]"]; ok {
				h(s, i, cfg)
				return
			}
			log.Warn("Failed to handle interaction", "ref", i.MessageComponentData().CustomID)
		case discordgo.InteractionModalSubmit:
			if h, ok := commandHandlers[strings.Split(
				i.ModalSubmitData().CustomID, "/")[0]+"[modal]"]; ok {
				h(s, i, cfg)
				return
			}
			log.Warn("Failed to handle interaction", "ref", i.ModalSubmitData().CustomID)
		default:
			log.Error("Unhandled interaction type", "type", i.Type)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: loc("generic.failed"),
			},
		})
	})

	err = s.Open()
	if err != nil {
		log.Fatal("Failed to open Discord session", "err", err)
	}
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, cfg config){
	loc("cmd.packs.title"): dMainCommand,
	"serverSelect[comp]":   dServerSelect,
	"switchTo[comp]":       dSwitchTo,
}
