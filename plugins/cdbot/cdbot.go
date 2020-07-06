package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/opcow/disgobot"
)

type bot string

var (
	// Bot is the exported bot
	Bot       bot
	lastCD    time.Time
	start     = make(chan int)
	quit      = make(chan bool)
	defaultCD = 5
	maxCD     = 5

	seed = rand.NewSource(time.Now().Unix())
	rnd  = rand.New(seed)

	idiot = []string{
		"That's not a number, idiot!",
		"Are you daft?",
		"I've done things with your mother.",
		"Your ass, your face, what's the difference?",
		"Did your mother have any children that lived?",
	}
)

func (b bot) BotInit(s []string) error {

	if len(s) > 1 {
		_, err := fmt.Sscanf(s[1], "%v", &defaultCD)
		if err != nil {
			fmt.Printf("cdbot: error parsing default CD: %s\n", err)
		} else {
			fmt.Printf("cdbot: default CD is %v.\n", defaultCD)
		}
	}
	if len(s) > 2 {
		_, err := fmt.Sscanf(s[2], "%v", &maxCD)
		if err != nil {
			fmt.Printf("cdbot: error parsing max CD: %s\n", err)
		} else {
			fmt.Printf("cdbot: max CD is %v.\n", maxCD)
		}
		if maxCD < 0 {
			maxCD = -maxCD
		}
		if defaultCD > maxCD {
			defaultCD = maxCD
		}
	}
	// Tell disgobot where to pass messages for processing
	return nil
}

func (b bot) BotExit() {
}

func (b bot) MessageProc(m *discordgo.MessageCreate, msg []string) bool {
	if msg[0] == "!cd" {
		if len(msg) > 1 {
			var count int
			_, err := fmt.Sscanf(msg[1], "%v", &count)
			//count, err := strconv.Atoi(m[1])
			if err != nil {
				disgobot.Discord.ChannelMessageSend(m.ChannelID, idiot[rnd.Intn(len(idiot))])
			} else if count == 0 {
				disgobot.Discord.ChannelMessageSend(m.ChannelID, "You in a hurry?")
			} else {
				if count < -maxCD || count > maxCD {
					disgobot.Discord.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The count must be a number from 1 to %v (defaults to %v).", maxCD, defaultCD))
				} else {
					printer(m.ChannelID, count)
				}
			}
		} else {
			printer(m.ChannelID, defaultCD)
		}
	}
	return true
}

// countdown printer
func printer(ChannelID string, n int) {
	if throttle(lastCD) {
		lastCD = time.Now()
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			if n == 0 {
				break
			}
			if n < 0 {
				disgobot.Discord.ChannelMessageSend(ChannelID, fmt.Sprintf("T-minus %v seconds...", -n))
				n++
			} else {
				disgobot.Discord.ChannelMessageSend(ChannelID, strconv.Itoa(n))
				n--
			}
		}
		ticker.Stop()
		disgobot.Discord.ChannelMessageSend(ChannelID, "_Go!_")
	} else {
		disgobot.Discord.ChannelMessageSend(ChannelID, "_No!_")
	}
}

// throttles responses
func throttle(lastTime time.Time) bool {
	//t1 := time.Date(2006, 1, 1, 12, 23, 10, 0, time.UTC)
	if time.Now().Sub(lastTime).Seconds() < 20 {
		return false
	}
	return true
}
