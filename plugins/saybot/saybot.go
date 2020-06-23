package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/opcow/disgobot"
)

type bot string

var (
	// Bot is the exported bot
	Bot           bot
	messageProcID string
	target        string
)

func (b bot) BotInit(s []string) error {
	fmt.Println("Saybot loaded.")
	return nil
}

func (b bot) BotExit() {
}

func (b bot) MessageProc(m *discordgo.MessageCreate, msg []string) bool {
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
	return false
}

func isChan(s string) bool {
	if _, err := disgobot.Discord.Channel(s); err == nil {
		return true
	}
	return false
}

// func listener() {
// 	var handler telnet.Handler = SayHandler
// 	err := telnet.ListenAndServe(":5555", handler)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// var SayHandler telnet.Handler = myHandler{}

// type myHandler struct{}

// func (handler myHandler) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {

// 	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
// 	p := buffer[:]

// 	for {
// 		n, err := r.Read(p)

// 		if n > 0 {
// 			oi.LongWrite(w, p[:n])
// 		}

// 		if nil != err {
// 			break
// 		}
// 	}
// }
