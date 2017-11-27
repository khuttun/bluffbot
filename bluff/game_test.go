package bluff

import (
	"testing"
)

func TestAddPlayer(t *testing.T) {
	var g Game
	a := PlayerInfo{42, "Alice"}
	b := PlayerInfo{43, "Bob"}
	e1 := g.AddPlayer(a)
	e2 := g.AddPlayer(b)
	if e1 != nil {
		t.Fail()
	}
	if e2 != nil {
		t.Fail()
	}
	if len(g.Players) != 2 {
		t.Fail()
	}
	if g.Players[0].Info != a {
		t.Fail()
	}
	if g.Players[1].Info != b {
		t.Fail()
	}
}

func TestAddPlayerWrongState(t *testing.T) {
	var g Game
	g.State = STARTED
	e := g.AddPlayer(PlayerInfo{42, "Alice"})
	if e == nil {
		t.Fail()
	}
}

func TestAddPlayerTwice(t *testing.T) {
	var g Game
	a := PlayerInfo{42, "Alice"}
	g.AddPlayer(a)
	e := g.AddPlayer(a)
	if e == nil {
		t.Fail()
	}
	if len(g.Players) != 1 {
		t.Fail()
	}
	if g.Players[0].Info != a {
		t.Fail()
	}
}

func TestStartGame(t *testing.T) {
	var g Game
	g.AddPlayer(PlayerInfo{42, "Alice"})
	g.AddPlayer(PlayerInfo{43, "Bob"})
	e := g.StartGame()
	if e != nil {
		t.Fail()
	}
	if g.State != STARTED {
		t.Fail()
	}
	for _, p := range g.Players {
		if len(p.Hand) != N_DICE_PER_PLAYER {
			t.Fail()
		}
	}
	if g.TurnIdx != 0 {
		t.Fail()
	}
	if g.CurrentBid.Count != 0 {
		t.Fail()
	}
}

func TestStartGameAlreadyStarted(t *testing.T) {
	var g Game
	g.AddPlayer(PlayerInfo{42, "Alice"})
	g.AddPlayer(PlayerInfo{43, "Bob"})
	g.StartGame()
	e := g.StartGame()
	if e == nil {
		t.Fail()
	}
	if g.State != STARTED {
		t.Fail()
	}
}

func TestStartGameTooFewPlayers(t *testing.T) {
	var g Game
	g.AddPlayer(PlayerInfo{42, "Alice"})
	e := g.StartGame()
	if e == nil {
		t.Fail()
	}
	if g.State != NOT_STARTED {
		t.Fail()
	}
}

func TestBidding(t *testing.T) {
	var g Game
	a := PlayerInfo{42, "Alice"}
	b := PlayerInfo{43, "Bob"}
	g.AddPlayer(a)
	g.AddPlayer(b)
	g.StartGame()

	aBid := Bid{a.ID, THREE, 3}
	e := g.Bid(aBid)
	if e != nil {
		t.Fail()
	}
	if g.TurnIdx != 1 {
		t.Fail()
	}
	if g.CurrentBid != aBid {
		t.Fail()
	}

	bBid := Bid{b.ID, THREE, 4}
	e = g.Bid(bBid)
	if e != nil {
		t.Fail()
	}
	if g.TurnIdx != 0 {
		t.Fail()
	}
	if g.CurrentBid != bBid {
		t.Fail()
	}
}

var scores = []struct {
	score int
	bid Bid
}{
	{0, Bid{Count: 0, Dice: ONE}},
	{1, Bid{Count: 1, Dice: ONE}},
	{2, Bid{Count: 1, Dice: TWO}},
	{3, Bid{Count: 1, Dice: THREE}},
	{4, Bid{Count: 1, Dice: FOUR}},
	{5, Bid{Count: 1, Dice: FIVE}},
	{6, Bid{Count: 1, Dice: WILD}},
	{7, Bid{Count: 2, Dice: ONE}},
	{8, Bid{Count: 2, Dice: TWO}},
	{9, Bid{Count: 2, Dice: THREE}},
	{10, Bid{Count: 2, Dice: FOUR}},
	{11, Bid{Count: 2, Dice: FIVE}},
	{12, Bid{Count: 3, Dice: ONE}},
	{13, Bid{Count: 3, Dice: TWO}},
	{14, Bid{Count: 3, Dice: THREE}},
	{15, Bid{Count: 3, Dice: FOUR}},
	{16, Bid{Count: 3, Dice: FIVE}},
	{17, Bid{Count: 2, Dice: WILD}},
	{18, Bid{Count: 4, Dice: ONE}},
	{19, Bid{Count: 4, Dice: TWO}},
	{20, Bid{Count: 4, Dice: THREE}},
	{21, Bid{Count: 4, Dice: FOUR}},
	{22, Bid{Count: 4, Dice: FIVE}},
	{23, Bid{Count: 5, Dice: ONE}},
	{24, Bid{Count: 5, Dice: TWO}},
	{25, Bid{Count: 5, Dice: THREE}},
	{26, Bid{Count: 5, Dice: FOUR}},
	{27, Bid{Count: 5, Dice: FIVE}},
	{28, Bid{Count: 3, Dice: WILD}},
	{29, Bid{Count: 6, Dice: ONE}},
	{30, Bid{Count: 6, Dice: TWO}},
	{31, Bid{Count: 6, Dice: THREE}},
	{32, Bid{Count: 6, Dice: FOUR}},
	{33, Bid{Count: 6, Dice: FIVE}},
	{34, Bid{Count: 7, Dice: ONE}},
	{35, Bid{Count: 7, Dice: TWO}},
	{36, Bid{Count: 7, Dice: THREE}},
	{37, Bid{Count: 7, Dice: FOUR}},
	{38, Bid{Count: 7, Dice: FIVE}},
	{39, Bid{Count: 4, Dice: WILD}},
	{40, Bid{Count: 8, Dice: ONE}},
}

func TestBidScore(t *testing.T) {
	for _, score := range scores {
		s := score.bid.score()
		if s != score.score {
			t.Fail()
		}
	}
}

func TestBidFromScore(t *testing.T) {
	for _, score := range scores {
		b := bidFromScore(score.score)
		if b.Count != score.bid.Count || b.Dice != score.bid.Dice {
			t.Fail()
		}
	}
}

func TestBiddingGameNotStarted(t *testing.T) {
	var g Game
	a := PlayerInfo{42, "Alice"}
	b := PlayerInfo{43, "Bob"}
	g.AddPlayer(a)
	g.AddPlayer(b)

	e := g.Bid(Bid{a.ID, THREE, 3})
	if e == nil {
		t.Fail()
	}
}

func TestBiddingOutOfTurn(t *testing.T) {
	var g Game
	a := PlayerInfo{42, "Alice"}
	b := PlayerInfo{43, "Bob"}
	c := PlayerInfo{44, "Carl"}
	g.AddPlayer(a)
	g.AddPlayer(b)
	g.AddPlayer(c)
	g.StartGame()

	aBid := Bid{a.ID, THREE, 3}
	g.Bid(aBid)
	e := g.Bid(Bid{c.ID, FOUR, 3})
	if e == nil {
		t.Fail()
	}
	if g.TurnIdx != 1 {
		t.Fail()
	}
	if g.CurrentBid != aBid {
		t.Fail()
	}

	e = g.Bid(Bid{45, FOUR, 4})
	if e == nil {
		t.Fail()
	}
	if g.TurnIdx != 1 {
		t.Fail()
	}
	if g.CurrentBid != aBid {
		t.Fail()
	}
}

func TestZeroBid(t *testing.T) {
	var g Game
	a := PlayerInfo{42, "Alice"}
	b := PlayerInfo{43, "Bob"}
	g.AddPlayer(a)
	g.AddPlayer(b)
	g.StartGame()
	e := g.Bid(Bid{a.ID, WILD, 0})
	if e == nil {
		t.Fail()
	}
}

func TestBiddingTooLow(t *testing.T) {
	var g Game
	a := PlayerInfo{42, "Alice"}
	b := PlayerInfo{43, "Bob"}
	g.AddPlayer(a)
	g.AddPlayer(b)
	g.StartGame()

	aBid := Bid{a.ID, THREE, 3}
	g.Bid(aBid)

	e := g.Bid(Bid{b.ID, WILD, 1})
	if e == nil {
		t.Fail()
	}
	if g.TurnIdx != 1 {
		t.Fail()
	}
	if g.CurrentBid != aBid {
		t.Fail()
	}
}

func TestChallengingWhenBidderWinsRound(t *testing.T) {
	var g Game
	g.State = STARTED
	g.Players = []Player{
		Player{PlayerInfo{1, "A"}, []Dice{ONE, ONE, ONE}},
		Player{PlayerInfo{2, "B"}, []Dice{WILD, TWO, THREE}}}
	g.TurnIdx = 1
	g.CurrentBid = Bid{1, ONE, 2}
	r, e := g.ChallengeCurrentBid(2)
	if e != nil {
		t.Fail()
	}
	if r.Result != LOW_BID {
		t.Fail()
	}
	if r.LostDiceCount != 2 {
		t.Fail()
	}
	if g.State != STARTED {
		t.Fail()
	}
	if len(g.Players[0].Hand) != 3 {
		t.Fail()
	}
	if len(g.Players[1].Hand) != 1 {
		t.Fail()
	}
	if g.TurnIdx != 0 {
		t.Fail()
	}
	if g.CurrentBid.Count != 0 {
		t.Fail()
	}
}

func TestChallengingWhenChallengerWinsRound(t *testing.T) {
	var g Game
	g.State = STARTED
	g.Players = []Player{
		Player{PlayerInfo{1, "A"}, []Dice{WILD, FIVE, ONE}},
		Player{PlayerInfo{2, "B"}, []Dice{FOUR, TWO, THREE}}}
	g.TurnIdx = 1
	g.CurrentBid = Bid{1, FIVE, 3}
	r, e := g.ChallengeCurrentBid(2)
	if e != nil {
		t.Fail()
	}
	if r.Result != HIGH_BID {
		t.Fail()
	}
	if r.LostDiceCount != 1 {
		t.Fail()
	}
	if g.State != STARTED {
		t.Fail()
	}
	if len(g.Players[0].Hand) != 2 {
		t.Fail()
	}
	if len(g.Players[1].Hand) != 3 {
		t.Fail()
	}
	if g.TurnIdx != 1 {
		t.Fail()
	}
	if g.CurrentBid.Count != 0 {
		t.Fail()
	}
}

func TestChallengingWhenBidIsExactlyRight(t *testing.T) {
	var g Game
	g.State = STARTED
	g.Players = []Player{
		Player{PlayerInfo{101, "A"}, []Dice{WILD, FIVE, ONE, ONE, TWO}},
		Player{PlayerInfo{102, "B"}, []Dice{FOUR, FIVE, THREE, WILD, TWO}},
		Player{PlayerInfo{103, "C"}, []Dice{FIVE, TWO, ONE, TWO, WILD}}}
	g.TurnIdx = 2
	g.CurrentBid = Bid{102, WILD, 3}
	r, e := g.ChallengeCurrentBid(103)
	if e != nil {
		t.Fail()
	}
	if r.Result != EXACT_BID {
		t.Fail()
	}
	if r.LostDiceCount != 1 {
		t.Fail()
	}
	if g.State != STARTED {
		t.Fail()
	}
	if len(g.Players[0].Hand) != 4 {
		t.Fail()
	}
	if len(g.Players[1].Hand) != 5 {
		t.Fail()
	}
	if len(g.Players[2].Hand) != 4 {
		t.Fail()
	}
	if g.TurnIdx != 1 {
		t.Fail()
	}
	if g.CurrentBid.Count != 0 {
		t.Fail()
	}
}

func TestChallengingWhenGameEnds(t *testing.T) {
	var g Game
	g.State = STARTED
	g.Players = []Player{
		Player{PlayerInfo{0, "A"}, []Dice{FIVE}},
		Player{PlayerInfo{1, "B"}, []Dice{FOUR, FOUR, FOUR}}}
	g.TurnIdx = 0
	g.CurrentBid = Bid{1, FOUR, 1}
	r, e := g.ChallengeCurrentBid(0)
	if e != nil {
		t.Fail()
	}
	if r.Result != LOW_BID {
		t.Fail()
	}
	if r.LostDiceCount != 2 {
		t.Fail()
	}
	if g.State != FINISHED {
		t.Fail()
	}
	if len(g.Players[0].Hand) != 0 {
		t.Fail()
	}
	if len(g.Players[1].Hand) != 3 {
		t.Fail()
	}
}

func TestChallengingGameNotStarted(t *testing.T) {
	var g Game
	a := PlayerInfo{42, "Alice"}
	b := PlayerInfo{43, "Bob"}
	g.AddPlayer(a)
	g.AddPlayer(b)

	_, e := g.ChallengeCurrentBid(a.ID)
	if e == nil {
		t.Fail()
	}
}

func TestChallengingOutOfTurn(t *testing.T) {
	var g Game
	a := PlayerInfo{42, "Alice"}
	b := PlayerInfo{43, "Bob"}
	c := PlayerInfo{44, "Carl"}
	g.AddPlayer(a)
	g.AddPlayer(b)
	g.AddPlayer(c)
	g.StartGame()
	g.Bid(Bid{a.ID, FIVE, 2})
	_, e := g.ChallengeCurrentBid(c.ID)
	if e == nil {
		t.Fail()
	}
}

func TestChallengingWhenNoBid(t *testing.T) {
	var g Game
	a := PlayerInfo{42, "Alice"}
	b := PlayerInfo{43, "Bob"}
	g.AddPlayer(a)
	g.AddPlayer(b)
	g.StartGame()
	_, e := g.ChallengeCurrentBid(a.ID)
	if e == nil {
		t.Fail()
	}
}

func TestBidOrdering(t *testing.T) {
	// Same count, different dice
	if !isGreater(Bid{1, TWO, 1}, Bid{1, ONE, 1}) {
		t.Fail()
	}
	if isGreater(Bid{1, ONE, 1}, Bid{1, TWO, 1}) {
		t.Fail()
	}

	// Different count, same dice
	if !isGreater(Bid{1, FIVE, 3}, Bid{1, FIVE, 2}) {
		t.Fail()
	}
	if isGreater(Bid{1, FIVE, 2}, Bid{1, FIVE, 3}) {
		t.Fail()
	}

	// Different count, different dice
	if !isGreater(Bid{1, ONE, 3}, Bid{1, FOUR, 2}) {
		t.Fail()
	}
	if isGreater(Bid{1, FOUR, 2}, Bid{1, ONE, 3}) {
		t.Fail()
	}

	// Wild vs normal
	if !isGreater(Bid{1, WILD, 5}, Bid{1, FIVE, 9}) {
		t.Fail()
	}
	if isGreater(Bid{1, FIVE, 9}, Bid{1, WILD, 5}) {
		t.Fail()
	}

	// Normal vs wild
	if !isGreater(Bid{1, ONE, 6}, Bid{1, WILD, 3}) {
		t.Fail()
	}
	if isGreater(Bid{1, WILD, 3}, Bid{1, ONE, 6}) {
		t.Fail()
	}

	// Wild vs wild
	if !isGreater(Bid{1, WILD, 5}, Bid{1, WILD, 4}) {
		t.Fail()
	}
	if isGreater(Bid{1, WILD, 4}, Bid{1, WILD, 5}) {
		t.Fail()
	}

	// Equal bids
	if isGreater(Bid{1, THREE, 10}, Bid{1, THREE, 10}) {
		t.Fail()
	}
}

func TestFindingNextPlayer(t *testing.T) {
	p := []Player{
		Player{PlayerInfo{}, []Dice{FIVE, FIVE}},
		Player{PlayerInfo{}, []Dice{}},
		Player{PlayerInfo{}, []Dice{WILD, ONE}}}

	i, e := indexOfNextPlayerWithDice(p, 0)
	if i != 2 {
		t.Fail()
	}
	if e != nil {
		t.Fail()
	}

	i, e = indexOfNextPlayerWithDice(p, 1)
	if i != 2 {
		t.Fail()
	}
	if e != nil {
		t.Fail()
	}

	i, e = indexOfNextPlayerWithDice(p, 2)
	if i != 0 {
		t.Fail()
	}
	if e != nil {
		t.Fail()
	}

	// Test situations when next player can't be found
	p[2].Hand = nil
	_, e = indexOfNextPlayerWithDice(p, 0)
	if e == nil {
		t.Fail()
	}

	_, e = indexOfNextPlayerWithDice(nil, 0)
	if e == nil {
		t.Fail()
	}
}
