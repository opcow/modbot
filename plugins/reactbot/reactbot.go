package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/opcow/disgobot"
)

type search struct {
	Channels []string
	Text     string
	Reaction []string
	Regex    bool
}

type searches struct {
	Search []search
	React  bool
}

type reactSpec struct {
	Search   string
	Channels map[string]struct{}
	Reaction []string
	Regex    bool
}

type bot string

var (
	// Bot is the exported bot
	Bot    bot
	ctx    context.Context
	cancel context.CancelFunc

	onOrOff = map[bool]string{
		false: "off",
		true:  "on",
	}

	conf        searches
	reacts      []reactSpec
	doReactions bool
)

func readConfig(f string) {
	tomlData, err := ioutil.ReadFile(f) // just pass the file name
	if err == nil {
		if _, err := toml.Decode(string(tomlData), &conf); err == nil {
			reacts = nil
			doReactions = conf.React
			var spec reactSpec
			for _, r := range conf.Search {
				spec.Channels = make(map[string]struct{})
				for _, c := range r.Channels {
					spec.Channels[c] = struct{}{}
				}
				spec.Search = r.Text
				spec.Reaction = r.Reaction
				spec.Regex = r.Regex
				reacts = append(reacts, spec)
			}
		}
	}
}

// func writeConfig() {
// 	f, _ := os.Create("reactbot.cfg")
// 	w := bufio.NewWriter(f)
// 	fmt.Fprintf(w, "react = %t\n\n", doReactions)
// 	for k, v := range reacts {
// 		fmt.Fprint(w, "[[search]]\n")
// 		fmt.Fprintf(w, `text = "%s"\n`, k)
// 		fmt.Fprint(w, "reaction = [")
// 		for i, s := range v.Reaction {
// 			if i == 0 {
// 				fmt.Fprintf(w, `"%s"`, s)
// 			} else {
// 				fmt.Fprintf(w, `, "%s"`, s)
// 			}
// 		}
// 		fmt.Fprint(w, "]\n\n")
// 	}
// 	w.Flush()
// }

func (b bot) BotInit(s []string) error {
	readConfig(path.Join(path.Dir(s[0]), "reactbot.cfg"))
	return nil
}

func (b bot) BotExit() {
	fmt.Println("ReactBot exiting...")
}

func (b bot) MessageProc(m *discordgo.MessageCreate, msg []string) bool {

	if msg[0] == "!react" {
		if len(msg) > 1 {
			switch msg[1] {
			case "on":
				doReactions = true
			case "off":
				doReactions = false
				// case "reload":
				// 	if disgobot.IsOp(m.Author.ID) {
				// 		readConfig()
				// 	}
			}
		}
		disgobot.Discord.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Reactions are %s", onOrOff[doReactions]))
	}

	if doReactions {
		var match bool
		for _, v := range reacts {
			if len(v.Channels) != 0 {
				if _, ok := v.Channels[m.ChannelID]; !ok {
					continue
				}
			}
			if v.Regex {
				match, _ = regexp.MatchString(v.Search, m.Content)
			} else {
				match = strings.Contains(strings.ToLower(m.Content), v.Search)
			}
			if match {
				for _, r := range v.Reaction {
					disgobot.Discord.MessageReactionAdd(m.ChannelID, m.ID, r)
				}
			}
		}
	}
	return false
}
