package bluff

import (
	"fmt"
	"math/rand"
)

const N_DICE_PER_PLAYER = 5

type Dice int

const (
	WILD Dice = iota
	ONE
	TWO
	THREE
	FOUR
	FIVE
)

type GameState int

const (
	// The game has not been started, new players can still be added
	NOT_STARTED GameState = iota
	// The game is started, bids and challenges can be made
	STARTED
	// The game is finished, the player having > 0 dice left is the winner
	FINISHED
)

type PlayerInfo struct {
	ID   int
	Name string
}

type Player struct {
	Info PlayerInfo
	Hand []Dice
}

// Set new random dice for a player
func (p *Player) rollDice() {
	for i := range p.Hand {
		p.Hand[i] = Dice(rand.Intn(int(FIVE + 1)))
	}
}

// Remove n dice from a player
func (p *Player) lostDice(n int) {
	nToRemove := n
	if nToRemove > len(p.Hand) {
		nToRemove = len(p.Hand)
	}
	p.Hand = p.Hand[nToRemove:]
}

type Bid struct {
	PlayerID int
	Dice     Dice
	Count    int
}

// Calculate score from bid to enable comparing bids
func (b *Bid) score() int {
	if b.Count <= 0 {
		return 0
	} else if b.Dice == WILD {
		return 6 + (b.Count - 1) * 11
	}
	// Leave "space" for stars
	stars := b.Count / 2
	return (b.Count - 1) * 5 + int(b.Dice) + stars
}

// Return quotient and remainder
func divmod(x, y int) (quot, rem int) {
	return x / y, x % y
}

// Construct a bid from score value
func bidFromScore(score int) Bid {
	if score <= 0 {
		return Bid{Count: 0, Dice: ONE}
	}
	stars, starsRem := divmod(score + 5, 11)
	if starsRem == 0 {
		return Bid{Count: stars, Dice: WILD}
	}
	count, face := divmod(score - stars - 1, 5)
	return Bid{Count: count + 1, Dice: Dice(face + 1)}
}

type BidClass int

const (
	// The bid was lower than the dice count, challenger loses
	LOW_BID BidClass = iota
	// The bid was exactly the dice count, everyone except the bidder loses one dice
	EXACT_BID
	// The bid was too high, bidder loses
	HIGH_BID
)

type ChallengeResult struct {
	Result        BidClass
	LostDiceCount int
	ChallengedBid Bid
	Bidder        PlayerInfo
	Challenger    PlayerInfo
}

type Game struct {
	State      GameState
	Players    []Player
	TurnIdx    int
	CurrentBid Bid
}

type GameError struct {
	What string
}

func (e *GameError) Error() string {
	return e.What
}

func (g *Game) AddPlayer(p PlayerInfo) error {
	if g.State != NOT_STARTED {
		return &GameError{"Can't add players when the game has already started"}
	}
	for _, v := range g.Players {
		if v.Info.ID == p.ID {
			return &GameError{"Player already added"}
		}
	}

	g.Players = append(g.Players, Player{p, nil})
	return nil
}

func (g *Game) StartGame() error {
	if g.State != NOT_STARTED {
		return &GameError{"Game already started"}
	}
	if len(g.Players) < 2 {
		return &GameError{"At least two players are needed to play"}
	}

	g.State = STARTED
	for i := range g.Players {
		g.Players[i].Hand = make([]Dice, N_DICE_PER_PLAYER)
		g.Players[i].rollDice()
	}
	g.TurnIdx = 0
	g.CurrentBid = Bid{}
	return nil
}

func (g *Game) Bid(b Bid) error {
	if g.State != STARTED {
		return &GameError{"Game not started"}
	}
	if b.PlayerID != g.Players[g.TurnIdx].Info.ID {
		return &GameError{fmt.Sprintf("It's %v's turn", g.Players[g.TurnIdx].Info.Name)}
	}
	if b.Count < 1 {
		return &GameError{"You must bid at least 1 dice"}
	}
	if !isGreater(b, g.CurrentBid) {
		return &GameError{"You must make a higher bid than the current one"}
	}

	g.CurrentBid = b
	var e error
	g.TurnIdx, e = indexOfNextPlayerWithDice(g.Players, g.TurnIdx)
	if e != nil {
		panic(e)
	}

	return nil
}

func (g *Game) ChallengeCurrentBid(playerID int) (ChallengeResult, error) {
	if g.State != STARTED {
		return ChallengeResult{}, &GameError{"Game not started"}
	}
	if playerID != g.Players[g.TurnIdx].Info.ID {
		return ChallengeResult{}, &GameError{fmt.Sprintf("It's %v's turn", g.Players[g.TurnIdx].Info.Name)}
	}
	if g.CurrentBid.Count < 1 {
		return ChallengeResult{}, &GameError{"No bid has been made yet"}
	}

	actualCount := totalCount(g.Players, g.CurrentBid.Dice)
	bidderIdx := indexOfId(g.Players, g.CurrentBid.PlayerID)
	bidder := &g.Players[bidderIdx]
	challenger := &g.Players[g.TurnIdx]
	result := ChallengeResult{ChallengedBid: g.CurrentBid, Bidder: bidder.Info, Challenger: challenger.Info}

	switch {
	// Less dice found than the bid -> bidder loses, challenger starts next round
	case actualCount < g.CurrentBid.Count:
		nLost := g.CurrentBid.Count - actualCount
		bidder.lostDice(nLost)
		result.Result = HIGH_BID
		result.LostDiceCount = nLost

	// More dice found than the bid -> challenger loses, bidder starts the next round
	case actualCount > g.CurrentBid.Count:
		nLost := actualCount - g.CurrentBid.Count
		challenger.lostDice(nLost)
		g.TurnIdx = bidderIdx
		result.Result = LOW_BID
		result.LostDiceCount = nLost

	// Bid was exactly right -> everyone except bidder loses one dice, bidder starts next round
	default:
		for i := range g.Players {
			if i != bidderIdx {
				g.Players[i].lostDice(1)
			}
		}
		g.TurnIdx = bidderIdx
		result.Result = EXACT_BID
		result.LostDiceCount = 1
	}

	// Roll new hand for everyone
	for i := range g.Players {
		g.Players[i].rollDice()
	}

	g.CurrentBid = Bid{}

	// Check whether the game ended
	_, e := indexOfNextPlayerWithDice(g.Players, g.TurnIdx)
	if e != nil {
		g.State = FINISHED
	}

	return result, nil
}

// Compare two bids: is b1 > b2?
func isGreater(b1 Bid, b2 Bid) bool {
	return b1.score() > b2.score()
}

// Find index of next Player with > 0 dice
func indexOfNextPlayerWithDice(players []Player, currentIndex int) (int, error) {
	n := len(players)
	if n < 2 {
		return currentIndex, &GameError{"Couldn't find next player"}
	}

	i := (currentIndex + 1) % n
	for len(players[i].Hand) == 0 {
		i = (i + 1) % n
		if i == currentIndex {
			return i, &GameError{"Couldn't find next player"}
		}
	}
	return i, nil
}

// Get total count of a dice value among players
func totalCount(players []Player, value Dice) int {
	c := 0
	for _, p := range players {
		for _, d := range p.Hand {
			if d == value || d == WILD {
				c++
			}
		}
	}
	return c
}

func indexOfId(players []Player, id int) int {
	for i := range players {
		if players[i].Info.ID == id {
			return i
		}
	}
	return -1
}
