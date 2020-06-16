package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/opcow/disgobot"
)

type search struct {
	Text     string
	Reaction []string
	Regex    bool
}

type searches struct {
	Search []search
	React  bool
}

type reactSpec struct {
	Reaction []string
	Regex    bool
}

var conf searches

var (
	ctx    context.Context
	cancel context.CancelFunc

	onOrOff = map[bool]string{
		false: "off",
		true:  "on",
	}

	reacts      = make(map[string]reactSpec)
	doReactions bool
)

type bot string

var Bot bot

func readConfig() {
	tomlData, err := ioutil.ReadFile("plugins/reactbot.cfg") // just pass the file name
	if err == nil {
		if _, err := toml.Decode(string(tomlData), &conf); err == nil {
			doReactions = conf.React
			var spec reactSpec
			for _, r := range conf.Search {
				spec.Reaction = r.Reaction
				spec.Regex = r.Regex
				reacts[r.Text] = spec
			}
		}
	}
}

func writeConfig() {
	f, _ := os.Create("reactbot.cfg")
	w := bufio.NewWriter(f)
	fmt.Fprintf(w, "react = %t\n\n", doReactions)
	for k, v := range reacts {
		fmt.Fprint(w, "[[search]]\n")
		fmt.Fprintf(w, `text = "%s"\n`, k)
		fmt.Fprint(w, "reaction = [")
		for i, s := range v.Reaction {
			if i == 0 {
				fmt.Fprintf(w, `"%s"`, s)
			} else {
				fmt.Fprintf(w, `, "%s"`, s)
			}
		}
		fmt.Fprint(w, "]\n\n")
	}
	w.Flush()
}

func (b bot) BotInit(s []string) {
	readConfig()
	disgobot.AddMessageProc(messageCreate)
}

func (b bot) BotExit() {
	fmt.Println("ReactBot exiting...")
}

func messageCreate(m *discordgo.MessageCreate, msg []string) {

	if msg[0] == "!react" {
		if len(msg) > 1 {
			if msg[1] == "on" {
				doReactions = true
			}
			if msg[1] == "off" {
				doReactions = false
			}
		}
		disgobot.Discord.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Reactions are %s", onOrOff[doReactions]))
	}

	if doReactions {
		var match bool
		for k, v := range reacts {
			if v.Regex {
				match, _ = regexp.MatchString(k, m.Content)
			} else {
				match = strings.Contains(strings.ToLower(m.Content), k)
			}
			if match {
				for _, r := range v.Reaction {
					disgobot.Discord.MessageReactionAdd(m.ChannelID, m.ID, r)
				}
			}
		}
	}
}
