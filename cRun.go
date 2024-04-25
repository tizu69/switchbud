package main

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

type cRun struct {
	RemoveCommands bool `param:"" help:"Remove all commands by the bot from the server? Useful if you changed the command name"`
	FuckSecurity   bool `param:"" help:"Fixes 'failed to verify certificate' errors (not recommended)"`
}

func (l *cRun) Run() error {
	cfg := getConfig()

	if l.FuckSecurity {
		log.Warn("Disabling cert checks! This is not recommended!")
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	ccInit(cfg)
	setupDiscord(cfg, l.RemoveCommands)

	time.Sleep(3 * time.Second)
	log.Info("Press CTRL-C to exit")
	<-make(chan struct{})

	return nil
}
