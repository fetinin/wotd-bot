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

// Implements telebot.v2.Recipient interface.
func (c channel) Recipient() string {
	return c.id
}

// Interval that specifies how often we should check if it's time
// to send new word.
const notificationInterval = 15 * time.Minute
// If anythings fails, wait for this amount before retry.
const retryInterval = 5 * time.Minute

// Hour after which message should be sent.
var submitHour = convertToInt(os.Getenv("submitTime"))
// Telegram channel ID.
var notificationChannel = channel{id: os.Getenv("channelID")}
// Telegram bot token.
var botToken = os.Getenv("botToken")

// Convert string to integer
func convertToInt(string string) int {
	hour, err := strconv.ParseInt(string, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	return int(hour)
}

// Notify channel subscribers about new word of the day.
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

// Get a new word of the day.
func getMessageOfTheDay() (string, error) {
	word, err := wotd.GetWOTD()
	return fmt.Sprintln(word), err
}

// Runs infinitely and notifies users on daily basis.
func runDailyNotification(bot *tb.Bot) {
	var lastSubmitDay int
	for {
		currentTime := time.Now()
		if lastSubmitDay != currentTime.YearDay() && currentTime.Hour() >= submitHour {
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
		b.Send(m.Sender, time.Now().Format("Mon Jan 2 15:04:05"))
	})

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, "Hello, "+m.Sender.FirstName+"!")
	})

	go runDailyNotification(b)

	defer b.Stop()
	b.Start()
}
