package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/carlescere/scheduler"
	"github.com/opcow/disgobot"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type bot string

type calConf struct {
	CalID    string
	Channels []string
	Time     string
}

// Bot is the exported bot
var Bot bot

var (
	botName       = "calbot"
	messageProcID string
	config        *oauth2.Config
	calID         string
	conf          calConf
	pluginPath    string
	calChans      = make(map[string]struct{})
	// reportCron    *cron.Cron
	// cronSpec      = "0 0"
	reportTime = "00:00"
	bdCommands = map[string]struct{}{
		"!bd":        {},
		"!bday":      {},
		"!bidet":     {},
		"!birthday":  {},
		"!birthdays": {},
	}
)

// BotInit() receives args for the bot and returns any error
func (b bot) BotInit(s []string) error {

	pluginPath = path.Dir(s[0])

	if err := readConfig(path.Join(pluginPath, "calbot.cfg")); err != nil {
		return err
	}

	br, err := ioutil.ReadFile(path.Join(pluginPath, "credentials.json"))
	if err != nil {
		fmt.Printf("Unable to read client secret file: %v", err)
		return err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err = google.ConfigFromJSON(br, calendar.CalendarReadonlyScope)
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %v", err)
		return err
	}

	// reportCron = cron.New(cron.WithParser(cron.NewParser(cron.Minute | cron.Hour)))
	// _, err = reportCron.AddFunc(cronSpec, cronReport)
	scheduler.Every().Day().At(reportTime).Run(eventReport)
	if err != nil {
		fmt.Printf("%s: error: %s", botName, err)
		return err
	}
	// reportCron.Start()
	fmt.Printf("%s report time is %s\n", botName, reportTime)

	return nil
}

func (b bot) BotExit() {
}

func eventReport() {
	if len(calChans) > 0 {
		n, resp := getBirthdaysToday()
		if n > 0 {
			for c := range calChans {
				disgobot.Discord.ChannelMessageSend(c, "Happy Birthday")
				disgobot.Discord.ChannelMessageSend(c, resp)
			}
		}
	}
}

// messageProc() receives a doscordgo MessageCreate struct and the
// message content is split into an array of words
func (b bot) MessageProc(m *discordgo.MessageCreate, msg []string) bool {
	if _, ok := calChans[m.ChannelID]; ok {
		if _, ok := bdCommands[msg[0]]; ok {
			days := 30
			if len(msg) > 1 {
				_, err := fmt.Sscanf(msg[1], "%v", &days)
				if err != nil {
					days = 30
				} else if days < 0 {
					days = 0
				} else if days >= 364 {
					days = 365
				}
			}
			_, resp := getBirthdays(days)
			if resp != "" {
				disgobot.Discord.ChannelMessageSend(m.ChannelID, resp)
			}
			// } else if msg[0] == "!cal" {
			// 	disgobot.Discord.ChannelMessageSend(m.ChannelID, "https://calendar.google.com/calendar?cid=aWZ0aWhyYmM2OTZ2aXFobzU2bWc0MWs3ZDhAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ")
		}
	}
	return true
}

func readConfig(f string) error {
	tomlData, err := ioutil.ReadFile(f) // just pass the file name
	if err == nil {
		if _, err := toml.Decode(string(tomlData), &conf); err == nil {
			calID = conf.CalID
			// cronSpec = conf.Cronspec
			reportTime = conf.Time
			for _, c := range conf.Channels {
				calChans[c] = struct{}{}
			}

		}
	}
	return err
}

func getBirthdays(days int) (int, string) {

	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		return 0, fmt.Sprintf("Unable to retrieve Calendar client: %v", err)
	}

	now := time.Now()
	tmin := now.Format(time.RFC3339)
	var tmax string
	if days == 0 {
		tmax = now.Add(1 * time.Second).Format(time.RFC3339)
	} else {
		tmax = now.AddDate(0, 0, days).Format(time.RFC3339)
	}

	events, err := srv.Events.List(calID).ShowDeleted(false).
		SingleEvents(true).TimeMin(tmin).TimeMax(tmax).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		return 0, fmt.Sprintf("Unable to retrieve next ten birthdays: %v", err)
	}
	count := len(events.Items)
	if count == 0 {
		return 0, "No upcoming birthdays found."
	} else {
		var b bytes.Buffer
		b.WriteString("```")
		for _, item := range events.Items {
			date := item.Start.DateTime
			if date == "" {
				date = strings.Replace(item.Start.Date[5:], "-", "/", 1)
			}
			b.WriteString(fmt.Sprintf("%v (%v)\n", item.Summary, date))
		}
		b.WriteString("```")
		return count, b.String()
	}
}

func getBirthdaysToday() (int, string) {

	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		fmt.Printf("Unable to retrieve Calendar client: %v", err)
		return 0, ""
	}

	now := time.Now()
	tmin := now.Format(time.RFC3339)
	tmax := now.Add(1 * time.Second).Format(time.RFC3339)

	events, err := srv.Events.List(calID).ShowDeleted(false).
		SingleEvents(true).TimeMin(tmin).TimeMax(tmax).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		fmt.Printf("Unable to retrieve next ten birthdays: %v", err)
		return 0, ""
	}
	count := len(events.Items)
	if count == 0 {
		return 0, ""
	} else {
		var b bytes.Buffer
		b.WriteString("```")
		for _, item := range events.Items {
			b.WriteString(fmt.Sprintf("%v\n", item.Summary))
		}
		b.WriteString("```")
		return count, b.String()
	}
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := path.Join(pluginPath, "token.json")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
