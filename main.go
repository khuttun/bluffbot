package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/khuttun/bluffbot/bluff"
	"github.com/khuttun/bluffbot/telegram"
)

func main() {
	port := os.Getenv("PORT")
	username := os.Getenv("USERNAME")
	token := os.Getenv("TELEGRAM_TOKEN")
	webhook := os.Getenv("WEBHOOK")

	if port == "" || username == "" || token == "" || webhook == "" {
		fmt.Println("Expecting following environment variables to be set:")
		fmt.Println("PORT")
		fmt.Println("USERNAME")
		fmt.Println("TELEGRAM_TOKEN")
		fmt.Println("WEBHOOK")
		os.Exit(1)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	t := telegram.BotAPI{Port: port, TelegramURL: fmt.Sprintf("https://api.telegram.org/bot%v/", token)}
	b := bluff.NewBot(username, &t)
	t.UpdateHandler = b.HandleUpdate
	t.SetWebhook(webhook)
	t.StartReceivingUpdates()
}
