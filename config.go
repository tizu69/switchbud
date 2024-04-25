package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/zijiren233/yaml-comment"
	"gopkg.in/yaml.v3"
)

type config struct {
	Version string `yaml:"version"`

	CcUrl      string `yaml:"ccApi" hc:"The Crafty Controller API URL (probably https://localhost:8443/api)"`
	CcUsername string `yaml:"ccUsername" hc:"API credentials - create a user in the web panel!"`
	CcPassword string `yaml:"ccPassword"`

	DiscordToken string `yaml:"discordToken" hc:"Discord bot token"`
	GuildId      string `yaml:"guildId" hc:"Discord guild ID to register commands in, or leave empty to register globally"`
	EmojiGuildId string `yaml:"emojiGuildId" hc:"Discord guild ID to store emojis in"`

	LangKeys map[string]string `yaml:"langKeys" hc:"Strings that the bot displays to the user"`

	Servers       []configServer `yaml:"servers" hc:"Servers that the bot is allowed to manage"`
	ResourceSlots int            `yaml:"resourceSlots" hc:"Number of resource slots available across all servers"`
}

var defaultConfig = config{
	Version: "0.1.0",

	CcUrl:      "https://localhost:8443/api",
	CcUsername: "admin",
	CcPassword: "password",

	DiscordToken: "Bot <token>",
	GuildId:      "",
	EmojiGuildId: "",

	LangKeys: localizations,

	Servers: []configServer{
		{
			DisplayName:  "All The Mods 9",
			DisplayDesc:  "**ATM9** has over **400 mods** and countless quests and a built in proper endgame. Can you craft the **ATM Star**? Do you dare take on the **Gregstar**?",
			DisplayIcon:  "https://media.forgecdn.net/avatars/902/338/638350403793040080.png",
			DownloadLink: "https://www.curseforge.com/minecraft/modpacks/all-the-mods-9/files",

			ServerID:      "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			ResourceSlots: 8,
			BootTime:      120,
		},
	},
	ResourceSlots: 8,
}

func getConfig() config {
	var c config = defaultConfig

	if _, err := os.Stat("switchbud.yml"); err != nil {
		log.Fatal("Failed to read config, try running 'switchbud init'", "err", err)
	}

	raw, err := os.ReadFile("switchbud.yml")
	if err != nil {
		log.Fatal("Failed to read config", "err", err)
	}

	err = yaml.Unmarshal(raw, &c)
	if err != nil {
		log.Fatal("Failed to parse config", "err", err)
	}

	return c
}

func saveConfig(c config) error {
	raw, err := yamlcomment.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile("switchbud.yml", raw, 0644)
}
