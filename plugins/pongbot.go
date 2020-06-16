package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/opcow/disgobot"
)

type bot string

var Bot bot

func (b bot) BotInit(s []string) {
	// Tell disgobot where to pass messages for processing
	disgobot.AddMessageProc(messageProc)

	// print args to log
	for i, o := range s {
		fmt.Printf("#%d: %s\n", i, o)
	}
}

func (b bot) BotExit() {
}

// messageProc() receives a doscordgo MessageCreate struct and the
// message content is split into an array of words
func messageProc(m *discordgo.MessageCreate, msg []string) {
	if strings.ToLower(m.Content) == "ping" {
		disgobot.Discord.ChannelMessageSend(m.ChannelID, "PONG")
	}
}
