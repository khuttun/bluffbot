package telegram

// MsgSender defines an interface to send Telegram messages
type MsgSender interface {
	SendMessage(chatid int, text string)
}
