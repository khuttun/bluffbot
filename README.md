# bluffbot

bluffbot is a [Telegram bot](https://core.telegram.org/bots) for playing a game of [Liar's dice (aka Bluff)](https://en.wikipedia.org/wiki/Liar%27s_dice). It's implemented in Go. Build the project with

```
go build
```

When running the executable, you need to have the following environment variables set:

* PORT: The port the process starts listening for [updates from Telegram](https://core.telegram.org/bots/api#getting-updates)
* USERNAME: The bot's Telegram username
* TELEGRAM_TOKEN: Your [Telegram API token](https://core.telegram.org/bots/api#authorizing-your-bot)
* WEBHOOK: [Telegram webhook](https://core.telegram.org/bots/api#setwebhook), the address Telegram should send the updates intended for this bot

The bluffbot repo includes couple of different packages:

* bluff: Core game logic and the logic for the bot itself
* telegram: Functions and types used to interact with the Telegram API

The bluffbot repo includes the files needed to run the bot in [Heroku](https://www.heroku.com/home) (Procfile, vendor.json).
