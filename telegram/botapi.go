package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// BotAPI offers an interface to the Telegram bot API
type BotAPI struct {
	// Port to start listening for updates from Telegram
	Port string
	// Send Telegram API requests here
	TelegramURL string
	// Function to handle updates coming from Telegram bot API
	UpdateHandler func(Update)
}

// Set URL where Telegram bot API sends updates
func (b *BotAPI) SetWebhook(url string) {
	b.makeRequest("setWebhook", SetWebhookParams{url})
}

// Send Telegram message
func (b *BotAPI) SendMessage(chatid int, text string) {
	b.makeRequest("sendMessage", SendMessageParams{ChatID: chatid, Text: text})
}

// Send Telegram message and display custom keyboard for the users
func (b *BotAPI) SendMessageAndDisplayCustomKeyboard(chatid int, text string, kb [][]string) {
	keyb := make([][]KeyboardButton, len(kb))
	for row := range kb {
		keyb[row] = make([]KeyboardButton, len(kb[row]))
		for col := range kb[row] {
			keyb[row][col] = KeyboardButton{kb[row][col]}
		}
	}
	b.makeRequest("sendMessage", SendMessageParams{ChatID: chatid, Text: text, ReplyMarkup: &ReplyKeyboardMarkup{Keyboard: keyb}})
}

// Send Telegram message and remove current custom keyboard
func (b *BotAPI) SendMessageAndRemoveCustomKeyboard(chatid int, text string) {
	b.makeRequest("sendMessage", SendMessageParams{ChatID: chatid, Text: text, ReplyMarkup: &ReplyKeyboardRemove{RemoveKeyboard: true}})
}

// Start receiving updates from Telegram bot API. Blocks.
func (b *BotAPI) StartReceivingUpdates() {
	http.HandleFunc("/", b.httpReqHandler)
	fmt.Println("Starting to listen port", b.Port)
	http.ListenAndServe(fmt.Sprintf(":%v", b.Port), nil)
}

func (b *BotAPI) makeRequest(method string, params interface{}) {
	paramsJSONStr, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("makeRequest", method, string(paramsJSONStr))

	resp, err := http.Post(b.TelegramURL+method, "application/json", bytes.NewReader(paramsJSONStr))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	fmt.Println("API response", resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
}

func (b *BotAPI) httpReqHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-----------")
	fmt.Println("New request")
	fmt.Println(r.Method)
	fmt.Println(r.URL)

	decoder := json.NewDecoder(r.Body)
	var upd Update
	err := decoder.Decode(&upd)
	if err != nil {
		fmt.Println(err)
		return
	}

	updStr, _ := json.MarshalIndent(upd, "", "    ")
	fmt.Println(string(updStr))

	fmt.Println("-----------")

	if upd.Message == nil {
		fmt.Println("nil Message")
		return
	}

	if upd.Message.From == nil {
		fmt.Println("nil From")
		return
	}

	if upd.Message.Text == nil {
		fmt.Println("nil Text")
		return
	}

	if b.UpdateHandler == nil {
		fmt.Println("nil UpdateHandler")
	}

	b.UpdateHandler(upd)
}
