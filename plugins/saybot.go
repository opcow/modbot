package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/opcow/disgobot"
	"github.com/reiver/go-telnet"
)

type bot string

var (
	Bot           bot
	messageProcID string
	target        string
)

func (b bot) BotInit(s []string) {
	// Tell disgobot where to pass messages for processing
	messageProcID = disgobot.AddMessageProc(messageProc)
	fmt.Println("Saybot loaded.")
	go listener()
}

func (b bot) BotExit() {
}

func isChan(s string) bool {
	if _, err := disgobot.Discord.Channel(s); err == nil {
		return true
	}
	return false
}

// messageProc() receives a doscordgo MessageCreate struct and the
// message content is split into an array of words
func messageProc(m *discordgo.MessageCreate, msg []string) {
	if disgobot.IsOp(m.Author.ID) {
		switch strings.ToLower(msg[0]) {
		case "!say":
			if len(msg) >= 3 && isChan(msg[1]) {
				target = msg[1]
				disgobot.Discord.ChannelMessageSend(target, strings.Join(msg[2:], " "))
			}
		case "!rsay":
			if target != "" && isChan(target) {
				disgobot.Discord.ChannelMessageSend(target, strings.Join(msg[1:], " "))
			}
		}
	}
}

func listener() {
	var handler telnet.Handler = telnet.EchoHandler
	err := telnet.ListenAndServe(":5555", handler)
}
