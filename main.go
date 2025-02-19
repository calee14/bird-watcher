package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

type MissionInfo struct {
	datename string
	time     string
}

func send(missionData []MissionInfo) {
	godotenv.Load()

	pass := os.Getenv("EMAIL_APP_PASSWORD")
	from := os.Getenv("EMAIL_SOURCE")
	to := os.Getenv("EMAIL_TARGET")

	_, month, day := time.Now().Date()

	body := fmt.Sprintf("Today's date: %s %s\n\n\n", month.String(), strconv.Itoa(day))
	for i := 0; i < len(missionData); i++ {
		mission := missionData[i]
		body += mission.datename + "\n"
		body += mission.time + "\n\n"
	}

	// make mail message
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "sightings for "+month.String()+" "+strconv.Itoa(day))
	msg.SetBody("text/plain", body)

	// dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587, from, pass)
	dialer.TLSConfig = &tls.Config{ServerName: "smtp.gmail.com"}

	err := dialer.DialAndSend(msg)
	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Printf("successfully sent message on: %s %s", month.String(), strconv.Itoa(day))
}

func collectMissionData() []MissionInfo {
	var missions []MissionInfo
	var mu sync.Mutex

	collector := colly.NewCollector(
		colly.AllowedDomains("spaceflightnow.com"),
	)

	// find and parse data
	collector.OnHTML("div.datename + div.missiondata", func(e *colly.HTMLElement) {
		missionDate := e.DOM.Prev().Text()
		missionTime := e.Text
		parsedMissionDate := strings.ReplaceAll(strings.TrimSpace(missionDate), "\n", " ")
		parsedMissionTime := strings.ReplaceAll(strings.TrimSpace(missionTime), "\n", " ")

		mu.Lock()
		missions = append(missions, MissionInfo{
			datename: parsedMissionDate,
			time:     parsedMissionTime,
		})
		mu.Unlock()

		// fmt.Println(parsedMissionDate)
		// fmt.Println(parsedMissionTime)
	})

	err := collector.Visit("https://spaceflightnow.com/launch-schedule/")
	if err != nil {
		log.Printf("error visiting page: %v", err)
	}
	return missions
}

func handleCli() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("exit", text) == 0 {
			log.Println("bird-watcher going offline.")
			os.Exit(0)
		}
	}
}

func watcher() {
	// collect data and send mail message
	missionData := collectMissionData()
	send(missionData)
}

func main() {
	go handleCli()

	loc, _ := time.LoadLocation("America/Los_Angeles")
	job := cron.New(cron.WithLocation(loc))
	_, err := job.AddFunc("0 8 * * *", watcher)
	if err != nil {
		log.Fatal(err)
	}
	job.Start()
	defer job.Stop()

	select {}
}
