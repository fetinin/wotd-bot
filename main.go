package main

import (
	"time"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
	"github.com/art-vasilyev/wotd"
	"fmt"
	"os"
)

type channel struct {
	id string
}

func (c channel) Recipient () string {
	return c.id
}

var notificationInterval = 10 * time.Minute
var submitHour = 9
var notificationChannel = channel{id: os.Getenv("channelID")}
var botToken = os.Getenv("botToken")


func notifySubscribers(b *tb.Bot) (err error) {
	log.Println("Notifying subscribers...")
	msg, err := getMessageOfTheDay()
	if err != nil {
		return
	}
	_, err = b.Send(notificationChannel, msg)
	if err != nil {
		return
	}
	return
}

func getMessageOfTheDay() (string, error){
	word, err := wotd.GetWOTD()
	return fmt.Sprintln(word), err
}

func runDailyNotification(bot *tb.Bot) {
	lastSubmitDay := 0
	for {
		if lastSubmitDay < time.Now().Day() {
			if time.Now().Hour() >= submitHour {
				err := notifySubscribers(bot)
				if err != nil {
					log.Println(err)
					continue // try again
				}
				lastSubmitDay = time.Now().Day()
			}
		}
		time.Sleep(notificationInterval)
	}
}

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	log.Printf("%s (@%v) started doing it's job.", b.Me.FirstName, b.Me.Username)
	log.Printf("Messages will be sent every day after %d:00 AM.", submitHour)

	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/ping", func(m *tb.Message) {
		b.Send(m.Sender, "♥️")
	})

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "Hello, " + m.Sender.FirstName +"!")
	})

	go runDailyNotification(b)

	defer b.Stop()
	b.Start()
}
