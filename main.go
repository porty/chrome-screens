package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/porty/chrome-screens/bot"
	"github.com/porty/chrome-screens/chrome"
	"github.com/porty/chrome-screens/config"

	"github.com/nlopes/slack"
)

type Settings struct {
	LastURL string `json:lastUrl`
}

func main() {
	mu := sync.Mutex{}
	cfg := config.Load()

	if err := cfg.Save(); err != nil {
		log.Print("Failed to save config: " + err.Error())
	}

	c, err := chrome.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialise Chrome wrapper: "+err.Error())
		os.Exit(1)
	}
	if cfg.LastURL != "" {
		c.SetURL(cfg.LastURL)
	}

	onURLUpdate := func(newURL string) {
		mu.Lock()
		defer mu.Unlock()

		cfg.LastURL = newURL
		if err := cfg.Save(); err != nil {
			log.Print("Error saving config file: " + err.Error())
		}
	}

	api := slack.New(cfg.SlackAPIToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	b := bot.New(rtm, c, onURLUpdate)
	b.Run()
}
