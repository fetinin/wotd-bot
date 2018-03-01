package main

import (
	"log"
	"time"

	"fmt"
	"github.com/art-vasilyev/wotd"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
	"strconv"
)

type channel struct {
	id string
}

func (c channel) Recipient() string {
	return c.id
}

var notificationInterval = 10 * time.Minute
var retryInterval = 5 * time.Minute
var submitHour = convertToInt(os.Getenv("submitTime"))
var notificationChannel = channel{id: os.Getenv("channelID")}
var botToken = os.Getenv("botToken")

func convertToInt(string string) int {
	hour, err := strconv.ParseInt(string, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	return int(hour)
}

func notifySubscribers(b *tb.Bot) (err error) {
	log.Println("Notifying subscribers...")
	msg, err := getMessageOfTheDay()
	if err != nil {
		return err
	}
	_, err = b.Send(notificationChannel, msg)
	if err != nil {
		return err
	}
	return err
}

func getMessageOfTheDay() (string, error) {
	word, err := wotd.GetWOTD()
	return fmt.Sprintln(word), err
}

func runDailyNotification(bot *tb.Bot) {
	lastSubmitDay := 0
	for {
		currentTime := time.Now()
		if lastSubmitDay < currentTime.Day() && currentTime.Hour() >= submitHour {
			err := notifySubscribers(bot)
			if err != nil {
				log.Println(err)
				log.Println("Will retry after", retryInterval)
				time.Sleep(retryInterval)
				continue // try again
			}
			lastSubmitDay = time.Now().Day()
		}
		time.Sleep(notificationInterval)
	}
}

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s (@%v) started doing it's job.", b.Me.FirstName, b.Me.Username)
	log.Printf("Messages will be sent every day after %d:00 AM.", submitHour)

	b.Handle("/ping", func(m *tb.Message) {
		b.Send(m.Sender, "♥️")
	})

	b.Handle("/time", func(m *tb.Message) {
		b.Send(m.Sender, time.Now())
	})

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "Hello, "+m.Sender.FirstName+"!")
	})

	go runDailyNotification(b)

	defer b.Stop()
	b.Start()
}
