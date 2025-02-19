package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
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

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: sightings for " + month.String() + " " + strconv.Itoa(day) + "\n\n" +
		body
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	log.Printf("Successfully sent message on: %s %s", month.String(), strconv.Itoa(day))
}

func collectMissionData() []MissionInfo {
	var missions []MissionInfo
	var mu sync.Mutex

	c := colly.NewCollector(
		colly.AllowedDomains("spaceflightnow.com"),
	)

	// find and parse data
	c.OnHTML("div.datename + div.missiondata", func(e *colly.HTMLElement) {
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

		fmt.Println(parsedMissionDate)
		fmt.Println((parsedMissionTime))
	})

	err := c.Visit("https://spaceflightnow.com/launch-schedule/")
	if err != nil {
		log.Printf("error visiting page: %v", err)
	}
	return missions
}

func main() {
	missionData := collectMissionData()
	send(missionData)
}
