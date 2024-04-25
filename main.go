package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"

	_ "github.com/bwmarrin/discordgo"
)

var cli struct {
	Init cInit `cmd:""`
	Run  cRun  `cmd:""`
}

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)

	ctx := kong.Parse(&cli)
	if _, err := os.Stat("switchbud.yml"); err == nil && ctx.Command() != "init" {
		cfg := getConfig()
		if cfg.Version != defaultConfig.Version {
			log.Fatal("Thanks for updating SwitchBud! Please run 'switchbud init' first to upgrade your config!")
		}
	}

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
