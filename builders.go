package main

import (
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

func bSingleComp(c discordgo.MessageComponent) *[]discordgo.MessageComponent {
	return &[]discordgo.MessageComponent{
		&discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				c,
			},
		},
	}
}

func bButton(label string, style discordgo.ButtonStyle, customId string, disabled bool) discordgo.MessageComponent {
	return discordgo.Button{
		Label:    label,
		Style:    style,
		CustomID: customId,
		Disabled: disabled,
	}
}

func bEmojiFromUrl(url, guildId string) *discordgo.ComponentEmoji {
	image, err := http.Get(url)
	if err != nil {
		log.Error("[imgoji] Failed to download image", "url", url, "err", err)
		return &discordgo.ComponentEmoji{Name: "❌"}
	}

	defer image.Body.Close()
	b, _ := io.ReadAll(image.Body)

	emoji, err := s.GuildEmojiCreate(guildId, &discordgo.EmojiParams{
		Image: "data:image/png;base64," + base64.StdEncoding.EncodeToString(b),
		Name:  "emoj",
	})
	if err != nil {
		log.Error("[imgoji] Failed to upload image", "url", url, "err", err)
		return &discordgo.ComponentEmoji{Name: "❌"}
	}

	time.AfterFunc(10*time.Second, func() {
		_ = s.GuildEmojiDelete(guildId, emoji.ID)
	})

	return &discordgo.ComponentEmoji{ID: emoji.ID, Name: emoji.Name}
}
