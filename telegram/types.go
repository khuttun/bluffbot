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
