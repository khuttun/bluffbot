package bluff

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/khuttun/bluffbot/telegram"
)

type Bot struct {
	username string
	telegram telegram.MsgSender
	games    map[int]*Game
}

// Create a new bot
func NewBot(uname string, tgram telegram.MsgSender) *Bot {
	return &Bot{username: uname, telegram: tgram, games: make(map[int]*Game)}
}

// Handle update from Telegram
func (b *Bot) HandleUpdate(u telegram.Update) {
	cmdParts := strings.Split(*u.Message.Text, " ")
	cmdName := strings.TrimSuffix(cmdParts[0], "@"+b.username)
	switch cmdName {
	case startCmd:
		b.onStartCmd(cmdParts[1:], *u.Message)
	case stopCmd:
		b.onStopCmd(cmdParts[1:], *u.Message)
	case beginCmd:
		b.onBeginCmd(cmdParts[1:], *u.Message)
	case bidCmd, bidButtonText:
		b.onBidCmd(cmdParts[1:], *u.Message)
	case challengeCmd:
		b.onChallengeCmd(cmdParts[1:], *u.Message)
	default:
		b.telegram.SendMessage(u.Message.Chat.ID, fmt.Sprintf("Unknown command: %v", cmdName))
	}
}

func (b *Bot) onStartCmd(params []string, msg telegram.Message) {
	if len(params) == 0 { // No params: start game
		if _, gameFound := b.games[msg.Chat.ID]; !gameFound {
			response := fmt.Sprintf("Starting a new game of %v! ", gameName)
			response += "Use the below link and click the START button in the opened chat window to join the game. "
			response += fmt.Sprintf("Once everyone has joined, send %v command to begin the game.", beginCmd)
			response += "\n\n"
			response += fmt.Sprintf("https://telegram.me/%v?start=%v", b.username, msg.Chat.ID)
			b.startGame(msg.Chat.ID, response)
		} else {
			b.telegram.SendMessage(msg.Chat.ID, "There's already a game started in this chat")
		}
	} else { // 1st param is the game/chat ID: player joining
		gameid, convErr := strconv.Atoi(params[0])
		if convErr != nil {
			b.telegram.SendMessage(msg.Chat.ID, fmt.Sprintf("Invalid game ID: %v", params[0]))
			return
		}

		if g, gameFound := b.games[gameid]; gameFound {
			err := g.AddPlayer(PlayerInfo{ID: msg.From.ID, Name: msg.From.FirstName})
			if err == nil {
				b.telegram.SendMessage(gameid, fmt.Sprintf("%v joined", msg.From.FirstName))
			} else {
				b.telegram.SendMessage(msg.Chat.ID, err.Error())
			}
		} else {
			b.telegram.SendMessage(msg.Chat.ID, fmt.Sprintf("Invalid game ID: %v", gameid))
		}
	}
}

func (b *Bot) onStopCmd(params []string, msg telegram.Message) {
	if _, gameFound := b.games[msg.Chat.ID]; gameFound {
		b.finishGame(msg.Chat.ID, "Game ended")
	} else {
		b.telegram.SendMessage(msg.Chat.ID, "No game started in this chat")
	}
}

func (b *Bot) onBeginCmd(params []string, msg telegram.Message) {
	if g, gameFound := b.games[msg.Chat.ID]; gameFound {
		err := g.StartGame()
		if err == nil {
			response := "The game begins. All the players should have now received their first round hand from me as a private message."
			response += "\n\n"
			response += fmt.Sprintf("Send \"%v count dice\" command to make a bid.", bidCmd)
			response += fmt.Sprintf("Use \"*\" for wild. For example, to make a bid of five wilds, send command \"%v 5 *\".", bidCmd)
			response += "\n\n"
			response += fmt.Sprintf("Send %v command to challenge current bid.", challengeCmd)
			response += "\n\n"
			response += turnMsg(g)
			b.beginRound(&msg.Chat, g, response)
		} else {
			b.telegram.SendMessage(msg.Chat.ID, err.Error())
		}
	} else {
		b.telegram.SendMessage(msg.Chat.ID, "No game started in this chat")
	}
}

func (b *Bot) onBidCmd(params []string, msg telegram.Message) {
	g, gameFound := b.games[msg.Chat.ID]
	if !gameFound {
		b.telegram.SendMessage(msg.Chat.ID, "No game started in this chat")
		return
	}

	if len(params) != 2 {
		b.telegram.SendMessage(msg.Chat.ID, fmt.Sprintf("Send \"%v count dice\" command to make a bid.", bidCmd))
		return
	}

	count, errc := strconv.Atoi(params[0])
	if errc != nil {
		b.telegram.SendMessage(msg.Chat.ID, fmt.Sprintf("Invalid count: %v", params[0]))
		return
	}

	d, errd := stringToDice(params[1])
	if errd != nil {
		b.telegram.SendMessage(msg.Chat.ID, errd.Error())
		return
	}

	errBid := g.Bid(Bid{PlayerID: msg.From.ID, Dice: d, Count: count})
	if errBid != nil {
		b.telegram.SendMessage(msg.Chat.ID, errBid.Error())
		return
	}

	b.telegram.SendMessageAndDisplayCustomKeyboard(msg.Chat.ID, fmt.Sprintf("%v bid %v %vs.%v", msg.From.FirstName, count, diceToString(d), turnMsg(g)), keyboard(g))
}

func (b *Bot) onChallengeCmd(params []string, msg telegram.Message) {
	g, gameFound := b.games[msg.Chat.ID]
	if !gameFound {
		b.telegram.SendMessage(msg.Chat.ID, "No game started in this chat")
		return
	}

	// Collect player hands already before calling ChallengeCurrentBid, because the call rolls new dice for everyone
	response := ""
	for _, p := range g.Players {
		if len(p.Hand) > 0 {
			response += fmt.Sprintf("%v: %v\n", p.Info.Name, handToString(p.Hand))
		}
	}
	response += "\n"

	r, e := g.ChallengeCurrentBid(msg.From.ID)
	if e != nil {
		b.telegram.SendMessage(msg.Chat.ID, e.Error())
		return
	}

	winner := ""

	switch r.Result {
	case LOW_BID:
		response += fmt.Sprintf("%v's bid was good. %v loses %v dice.", r.Bidder.Name, r.Challenger.Name, r.LostDiceCount)
		winner = r.Bidder.Name
	case EXACT_BID:
		response += fmt.Sprintf("%v's bid was exactly right! Everyone else loses %v dice.", r.Bidder.Name, r.LostDiceCount)
		winner = r.Bidder.Name
	case HIGH_BID:
		response += fmt.Sprintf("%v's bid was too high. %v loses %v dice.", r.Bidder.Name, r.Bidder.Name, r.LostDiceCount)
		winner = r.Challenger.Name
	}

	response += "\n\n"
	response += gameStatusMsg(g)
	response += "\n\n"

	switch g.State {
	case STARTED:
		response += fmt.Sprintf("Starting next round. %v", turnMsg(g))
		b.beginRound(&msg.Chat, g, response)
	case FINISHED:
		response += fmt.Sprintf("Game finished! %v is the winner!", winner)
		b.finishGame(msg.Chat.ID, response)
	}
}

func (b *Bot) startGame(chatId int, msg string) {
	b.games[chatId] = &Game{}
	b.telegram.SendMessage(chatId, msg)
}

func (b *Bot) finishGame(chatId int, msg string) {
	b.telegram.SendMessageAndRemoveCustomKeyboard(chatId, msg)
	delete(b.games, chatId)
}

func (b *Bot) beginRound(chat *telegram.Chat, g *Game, msg string) {
	b.telegram.SendMessageAndDisplayCustomKeyboard(chat.ID, msg, keyboard(g))
	b.sendHands(g, chat.Title)
}

func (b *Bot) sendHands(g *Game, chatname *string) {
	var cn string
	if chatname != nil {
		cn = *chatname
	}
	for _, p := range g.Players {
		b.telegram.SendMessage(p.Info.ID, fmt.Sprintf("Your %v hand in %v:\n%v", gameName, cn, handToString(p.Hand)))
	}
}

const gameName = "Bluff"

const startCmd = "/start"
const stopCmd = "/stop"
const beginCmd = "/begin"
const bidCmd = "/bid"
const challengeCmd = "/challenge"
const bidButtonText = "Bid"

func gameStatusMsg(g *Game) string {
	msg := "Game status:"
	total := 0
	for _, p := range g.Players {
		msg += fmt.Sprintf("\n%v %v dice", p.Info.Name, len(p.Hand))
		total += len(p.Hand)
	}
	msg += fmt.Sprintf("\nTotal %v dice", total)
	return msg
}

func turnMsg(g *Game) string {
	return fmt.Sprintf("It's %v's turn.", g.Players[g.TurnIdx].Info.Name)
}

func diceToString(d Dice) string {
	switch d {
	case WILD:
		return "*️⃣"
	case ONE:
		return "1️⃣"
	case TWO:
		return "2️⃣"
	case THREE:
		return "3️⃣"
	case FOUR:
		return "4️⃣"
	case FIVE:
		return "5️⃣"
	}
	return "?"
}

func stringToDice(s string) (Dice, error) {
	switch s {
	case "*", diceToString(WILD):
		return WILD, nil
	case "1", diceToString(ONE):
		return ONE, nil
	case "2", diceToString(TWO):
		return TWO, nil
	case "3", diceToString(THREE):
		return THREE, nil
	case "4", diceToString(FOUR):
		return FOUR, nil
	case "5", diceToString(FIVE):
		return FIVE, nil
	default:
		return WILD, fmt.Errorf("Unknown dice: %v", s)
	}
}

func handToString(hand []Dice) string {
	s := ""
	for _, d := range hand {
		s += diceToString(d)
	}
	return s
}

func keyboard(g *Game) [][]string {
	b := g.CurrentBid
	kb := make([][]string, 4)
	for row := range kb {
		kb[row] = make([]string, 4)
		for col := range kb[row] {
			b = nextBid(b)
			kb[row][col] = bidButton(b)
		}
	}
	return kb
}

func bidButton(b Bid) string {
	return fmt.Sprintf("%v %v %v", bidButtonText, b.Count, diceToString(b.Dice))
}

func nextBid(b Bid) Bid {
	c := b.effectiveCount()
	if c < 1 {
		return Bid{Dice: ONE, Count: 1}
	}
	switch b.Dice {
	case WILD:
		return Bid{Dice: ONE, Count: c}
	case ONE:
		return Bid{Dice: TWO, Count: c}
	case TWO:
		return Bid{Dice: THREE, Count: c}
	case THREE:
		return Bid{Dice: FOUR, Count: c}
	case FOUR:
		return Bid{Dice: FIVE, Count: c}
	case FIVE:
		switch c % 2 {
		case 0:
			return Bid{Dice: ONE, Count: c + 1}
		case 1:
			return Bid{Dice: WILD, Count: c/2 + 1}
		}
	}
	return Bid{Dice: ONE, Count: 1}
}
