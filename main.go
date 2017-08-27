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
	args := os.Args[1:]

	if len(args) < 3 {
		fmt.Println("Usage:", os.Args[0], "bot_username telegram_token webhook_addr")
		os.Exit(1)
	}

	username := args[0]
	token := args[1]
	webhook := args[2]

	rand.Seed(time.Now().UTC().UnixNano())
	t := telegram.BotAPI{TelegramURL: fmt.Sprintf("https://api.telegram.org/bot%v/", token)}
	b := bluff.NewBot(username, &t)
	t.UpdateHandler = b.HandleUpdate
	t.SetWebhook(webhook)
	t.StartReceivingUpdates()
}
