package telegram

// User represents a Telegram user or bot.
type User struct {
	// Unique identifier for this user or bot.
	ID int `json:"id"`
	// User‘s or bot’s first name.
	FirstName string `json:"first_name"`
	// Optional. User‘s or bot’s last name.
	LastName *string `json:"last_name"`
	// Optional. User‘s or bot’s username.
	Username *string `json:"username"`
}

// Chat represents a Telegram chat.
type Chat struct {
	// Unique identifier for this chat.
	ID int `json:"id"`
	// Type of chat, can be either “private”, “group”, “supergroup” or “channel”.
	Type string `json:"type"`
	// Optional. Title, for supergroups, channels and group chats.
	Title *string `json:"title"`
}

// Message represents a Telegram message.
type Message struct {
	// Unique message identifier inside this chat.
	MessageID int `json:"message_id"`
	// Date the message was sent in Unix time.
	Date int `json:"date"`
	// Conversation the message belongs to.
	Chat Chat `json:"chat"`
	// Optional. Sender, can be empty for messages sent to channels.
	From *User `json:"from"`
	// Optional. For text messages, the actual UTF-8 text of the message, 0-4096 characters.
	Text *string `json:"text"`
}

// Update represents an incoming Telegram update.
type Update struct {
	// The update‘s unique identifier. Update identifiers start from a certain positive number and increase sequentially.
	UpdateID int `json:"update_id"`
	// Optional. New incoming message of any kind — text, photo, sticker, etc.
	Message *Message `json:"message"`
}

// SetWebhookParams defines parameters for Telegram API setWebhook method
type SetWebhookParams struct {
	// HTTPS url to send updates to. Use an empty string to remove webhook integration.
	URL string `json:"url"`
}

// KeyboardButton represents one button of the reply keyboard.
type KeyboardButton struct {
	// Text of the button. It will be sent to the bot as a message when the button is pressed.
	Text string `json:"text"`
}

// ReplyKeyboardMarkup represents a custom keyboard with reply options.
type ReplyKeyboardMarkup struct {
	// Array of button rows, each represented by an Array of KeyboardButton objects.
	Keyboard [][]KeyboardButton `json:"keyboard"`
	// Optional. Requests clients to resize the keyboard vertically for optimal fit.
	ResizeKeyboard *bool `json:"resize_keyboard,omitempty"`
	// Optional. Requests clients to hide the keyboard as soon as it's been used.
	OneTimeKeyboard *bool `json:"one_time_keyboard,omitempty"`
	// Optional. Use this parameter if you want to show the keyboard to specific users only.
	// Targets:
	// 1) users that are @mentioned in the text of the Message object
	// 2) if the bot's message is a reply (has reply_to_message_id), sender of the original message.
	Selective *bool `json:"selective,omitempty"`
}

// ReplyKeyboardRemove is used to remove custom keyboard from chat.
type ReplyKeyboardRemove struct {
	// Requests clients to remove the custom keyboard. Should always be true.
	RemoveKeyboard bool `json:"remove_keyboard"`
	// Optional. Use this parameter if you want to remove the keyboard for specific users only. Targets:
	// 1) users that are @mentioned in the text of the Message object;
	// 2) if the bot's message is a reply (has reply_to_message_id), sender of the original message.
	Selective *bool `json:"selective,omitempty"`
}

// SendMessageParams defines parameters for Telegram API sendMessage method
type SendMessageParams struct {
	// Unique identifier for the target chat
	ChatID int `json:"chat_id"`
	// Text of the message to be sent
	Text string `json:"text"`
	// Optional. Add/remove custom keyboard. Allowed types:
	// - ReplyKeyboardMarkup
	// - ReplyKeyboardRemove
	ReplyMarkup interface{} `json:"reply_markup,omitempty"`
}
