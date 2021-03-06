package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/opcow/disgobot"
)

type config struct {
	Token   string
	Plugins []string
	Ops     []string
}

var (
	conf     config
	confFile = flag.String("c", "", "config file")

	onOrOff = map[bool]string{
		false: "off",
		true:  "on",
	}
	addS = map[bool]string{
		false: "",
		true:  "s",
	}
	sc chan os.Signal
)

func readConfig() {
	var tomlData []byte
	var err error
	if *confFile != "" {
		tomlData, err = os.ReadFile(*confFile)
	} else {
		tomlData, err = os.ReadFile("modbot.cfg")
	}
	if err != nil {
		panic(err)
	}
	if _, err = toml.Decode(string(tomlData), &conf); err != nil {
		panic(err)
	}
}

func main() {

	flag.Parse()
	readConfig()

	if conf.Token == "" {
		fmt.Println("Auth token for bot not suppiled in config.")
		os.Exit(1)
	}

	for _, o := range conf.Ops {
		disgobot.AddOp(o)
	}

	err := disgobot.Run(conf.Token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, p := range conf.Plugins {
		err := disgobot.LoadPlugin(p)
		if err != nil {
			fmt.Printf("Error loading plugin %s: %s\n", p, err)
		}
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	signal.Notify(disgobot.SignalChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-disgobot.SignalChan

	// Cleanly close down the Discord session.
	disgobot.Discord.Close()
}
