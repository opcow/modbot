package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/opcow/disgobot"
	"github.com/robfig/cron/v3"
)

type tests struct {
	Total int `json:"total"`
}

type deaths struct {
	New   string `json:"new"`
	Total int    `json:"total"`
}

type cases struct {
	New       string `json:"new"`
	Active    int    `json:"active"`
	Critical  int    `json:"critical"`
	Recovered int    `json:"recovered"`
	Total     int    `json:"total"`
}

type response struct {
	Country string `json:"country"`
	Cases   cases  `json:"cases"`
	Deaths  deaths `json:"deaths"`
	Tests   tests  `json:"tests"`
	Day     string `json:"day"`
	Time    string `json:"time"`
}

type params struct {
	Country string `json:"country"`
}

type covidReport struct {
	Get        string     `json:"get"`
	Parameters params     `json:"parameters"`
	Errors     []int      `json:"errors"`
	Results    int        `json:"results"`
	Response   []response `json:"response"`
}

type config struct {
	Token    string
	Chans    []string
	Cronspec string
}

type bot string

var (
	// Bot is the exported bot
	Bot          bot
	botName      = "covidbot"
	conf         config
	seed         = rand.NewSource(time.Now().Unix())
	rnd          = rand.New(seed)
	token        string
	initialChans string
	cronSpec     = "1 0"
	// cronEntryID  cron.EntryID
	lastReport time.Time

	covChans   = make(map[string]struct{})
	reportCron *cron.Cron

	nfStrings = []string{
		"Must be a shithole.",
		"Perhaps they're all dead.",
		"How about a nice game of TIC-TAC-TOE?",
		"Thanks, Hillary.",
		"Try not spelling like the president.",
		"I felt a great disturbance in the force. Coincidence?",
		"Maybe it's fictional. Like Finland.",
		"Maybe you could discover it.",
	}
)

func readConfig(f string) error {
	fmt.Printf("%s: loading %s...\n", botName, f)
	tomlData, err := ioutil.ReadFile(f) // just pass the file name
	if err != nil {
		fmt.Printf("%s: error: %s", botName, err)
		return err
	}
	if _, err := toml.Decode(string(tomlData), &conf); err == nil {
		token = conf.Token
		cronSpec = conf.Cronspec
		for _, c := range conf.Chans {
			covChans[c] = struct{}{}
		}
	} else {
		fmt.Printf("%s: error: %s", botName, err)
	}
	return err
}

func (b bot) BotInit(s []string) error {
	var err error

	if err = readConfig(path.Join(path.Dir(s[0]), "covidbot.cfg")); err != nil {
		return err
	}
	if token == "" {
		fmt.Printf("%s: disabled due to empty token string.", botName)
		return errors.New("missing token")
	}
	reportCron = cron.New(cron.WithParser(cron.NewParser(cron.Minute | cron.Hour)))
	_, err = reportCron.AddFunc(cronSpec, cronReport)
	if err == nil {
		reportCron.Start()
	} else {
		fmt.Printf("%s: error: %s", botName, err)
	}
	fmt.Printf("Cronspec is %s\n", cronSpec)

	return nil
}

func (b bot) BotExit() {
	fmt.Println("ReactBot exiting...")
	// Stop cron jobs.
	reportCron.Stop()
}

func (b bot) MessageProc(m *discordgo.MessageCreate, msg []string) bool {
	switch msg[0] {
	case "!cov": // report covid-19 stats
		if time.Now().Sub(lastReport).Seconds() < 10 {
			disgobot.Discord.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Please wait %.0f seconds and try again.", 10.0-time.Now().Sub(lastReport).Seconds()))
			return false
		}
		var err error
		var report string
		// lastCD = time.Now()
		if len(msg) > 1 {
			report, err = covid(strings.Join(msg[1:], "-"))
		} else {
			report, err = covid("usa")
		}
		if err == nil {
			disgobot.Discord.ChannelMessageSend(m.ChannelID, report)
		}
	case "!reaper": // periodic USA death toll reports
		if len(msg) < 2 || msg[1] != "off" {
			if !disgobot.IsOp(m.Author.ID) {
				return false
			}
			if len(msg) == 1 {
				// just report the status
				disgobot.Discord.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Grim Reaper reports are *on* for %s.", disgobot.ChanIDtoMention(m.ChannelID)))
				covChans[m.ChannelID] = struct{}{}
			} else if id, err := disgobot.ChanMentionToID(msg[1]); err == nil {
				covChans[id] = struct{}{}
				disgobot.Discord.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Grim Reaper reports are *on* for %s.", disgobot.ChanIDtoMention(m.ChannelID)))
			}
		} else if len(msg) == 2 {
			disgobot.Discord.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Grim Reaper reports are *off* for %s.", disgobot.ChanIDtoMention(m.ChannelID)))
			delete(covChans, m.ChannelID)
		} else if id, err := disgobot.ChanMentionToID(msg[2]); err == nil {
			delete(covChans, id)
			disgobot.Discord.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Grim Reaper reports are *off* for %s.", disgobot.ChanIDtoMention(m.ChannelID)))
		}
	case "!covchans":
		if !disgobot.IsOp(m.Author.ID) {
			return false
		}
		c, err := disgobot.Discord.UserChannelCreate(m.Author.ID)
		if err != nil {
			return false
		}
		var s = "channels:"
		for k := range covChans {
			s = s + " " + disgobot.ChanIDtoMention(k)
		}
		disgobot.Discord.ChannelMessageSend(c.ID, s)
		time.Sleep(time.Millisecond * 500)
		disgobot.Discord.ChannelMessageSend(c.ID, fmt.Sprintf("cronspec: %s", cronSpec))
	}
	return false
}

func covid(country string) (string, error) {

	var report covidReport
	var newDeaths string

	url := "https://covid-193.p.rapidapi.com/statistics?country=" + country
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-host", "covid-193.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", token)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &report)

	if report.Results < 1 {
		return fmt.Sprintf("No results for %s. %s", country, nfStrings[rnd.Intn(len(nfStrings))]), nil
	}

	if report.Response[0].Deaths.New != "" {
		newDeaths = fmt.Sprintf(". (%s)", report.Response[0].Deaths.New)
	} else {
		newDeaths = "."
	}

	if country == "all" {
		return fmt.Sprintf("Covid-19 World: %d active cases, %d critical cases, %d recoverd, %d total cases, %d deaths%s\n",
			report.Response[0].Cases.Active, report.Response[0].Cases.Critical, report.Response[0].Cases.Recovered,
			report.Response[0].Cases.Total, report.Response[0].Deaths.Total, newDeaths), nil
	}
	return fmt.Sprintf("Covid-19 %s: %d tested, %d active cases, %d critical cases, %d recoverd, %d total cases, %d deaths%s\n",
		report.Response[0].Country, report.Response[0].Tests.Total, report.Response[0].Cases.Active, report.Response[0].Cases.Critical,
		report.Response[0].Cases.Recovered, report.Response[0].Cases.Total, report.Response[0].Deaths.Total, newDeaths), nil
}

func reaper() (string, error) {

	var report covidReport

	url := "https://covid-193.p.rapidapi.com/statistics?country=usa"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-host", "covid-193.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", token)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &report)

	if report.Results < 1 {
		return "No death count available.", nil
	}

	t, _ := time.Parse(time.RFC3339, report.Response[0].Time)
	location, err := time.LoadLocation("America/New_York")
	var tStr string
	// var tLoc time.Time

	if err != nil {
		tStr = report.Response[0].Time
	} else {
		tLoc := t.In(location)
		zone, _ := tLoc.Zone()
		tStr = tLoc.Format("2006-01-02 @ 15:04 ") + zone
	}

	var newDeaths string
	if report.Response[0].Deaths.New != "" {
		newDeaths = fmt.Sprintf(". (%s)", report.Response[0].Deaths.New)
	} else {
		newDeaths = "."
	}

	return fmt.Sprintf("USA (%s): %d covid-19 deaths%s\n", tStr, report.Response[0].Deaths.Total, newDeaths), nil
}

func cronReport() {
	if len(covChans) > 0 {
		report, err := reaper()
		if err == nil {
			for c := range covChans {
				disgobot.Discord.ChannelMessageSend(c, report)
			}
		}
	}
}
