package watcher

import (
	db "bird-watcher/internal/database"
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
)

type MissionInfo struct {
	datename string
	time     string
}

type ScheduledDate struct {
	Day   int
	Month int
}

// scheduled time to send updates
var ScheduledTime int = 8

// init prev scheduled date to trigger send
var PrevScheduledDate ScheduledDate = ScheduledDate{time.Now().Day() - 1, int(time.Now().Month())}

func Send(missionData []MissionInfo, target string) {
	godotenv.Load()

	pass := os.Getenv("EMAIL_APP_PASSWORD")
	from := os.Getenv("EMAIL_SOURCE")
	to := target

	_, month, day := time.Now().Date()

	subject := fmt.Sprintf("%s %s | Upcoming: %s", month.String(), strconv.Itoa(day), missionData[0].datename)
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
	msg.SetHeader("Subject", subject)
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

func CollectMissionData() []MissionInfo {
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

func HandleCli() {
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

func Watcher() {
	// collect data and send mail message
	missionData := CollectMissionData()

	subscribers, err := db.GetAllSubscribers()
	if err != nil {
		log.Printf("error: %s", err.Error())
	}

	// send out emails in goroutine
	var wg sync.WaitGroup
	for i := 0; i < len(subscribers); i++ {
		wg.Add(1)

		go func(data interface{}, email string) {
			defer wg.Done()

			Send(missionData, email)
		}(missionData, subscribers[i].Email)
	}
	wg.Wait()
}

func StartWatcher() {
	for {
		loc, _ := time.LoadLocation("America/Los_Angeles")
		now := time.Now().In(loc)
		_, mm, dd := now.Date()
		hour := now.Hour()
		currDate := ScheduledDate{dd, int(mm)}
		if hour == ScheduledTime && PrevScheduledDate != currDate {
			Watcher()
			PrevScheduledDate = ScheduledDate{dd, int(mm)}
		}
		time.Sleep(5 * time.Minute)
	}
}

// func StartWatcher() {
// 	loc, _ := time.LoadLocation("America/Los_Angeles")
// 	job := cron.New(cron.WithLocation(loc))
// 	_, err := job.AddFunc("* * * * *", Watcher)
// 	if err != nil {
// 		log.Fatal(err)
// 		log.Printf("starting cron job. errors: %s", err.Error())
// 	}
// 	job.Start()
// 	defer job.Stop()
// }
