package telegram

// MsgSender defines an interface to send Telegram messages
type MsgSender interface {
	SendMessage(chatid int, text string)
	SendMessageAndDisplayCustomKeyboard(chatid int, text string, kb [][]string)
	SendMessageAndRemoveCustomKeyboard(chatid int, text string)
}
