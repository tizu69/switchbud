package main

import (
	"os"

	"github.com/charmbracelet/log"
)

type cInit struct {
}

func (l *cInit) Run() error {
	if _, err := os.Stat("switchbud.yml"); err == nil {
		cfg := getConfig()
		cfg.Version = defaultConfig.Version
		saveConfig(cfg)

		log.Info("Cleaned up config -- formatted, added missing options and fixed comments!")
		return nil
	}

	err := saveConfig(defaultConfig)
	if err != nil {
		log.Fatal("Oops! Failed to save config", "err", err)
		return err
	}

	log.Info("Open ./switchbud.yml in your favorite editor to configure!")
	return nil
}
