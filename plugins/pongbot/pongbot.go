package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/opcow/disgobot"
)

type bot string

// Bot is the exported bot
var Bot bot
var messageProcID string

// BotInit() receives args for the bot and returns any error
func (b bot) BotInit(s []string) error {

	// print args to log
	for i, o := range s {
		fmt.Printf("#%d: %s\n", i, o)
	}
	return nil
}

func (b bot) BotExit() {
}

// messageProc() receives a doscordgo MessageCreate struct and the
// message content is split into an array of words
func (b bot) MessageProc(m *discordgo.MessageCreate, msg []string) bool {
	switch strings.ToLower(m.Content) {
	case "ping":
		disgobot.Discord.ChannelMessageSend(m.ChannelID, "PONG")
	case "nopongs":
		disgobot.Discord.ChannelMessageSend(m.ChannelID, "No PONGS for you.")
	}
	return true
}
