// TODO(refactor): Change return types to (<type>, error) and handle errors
package blackjack

import (
	"github.com/GRO4T/blackjack/deck"
)

const (
	initialChips = 100
)

type State int

const (
	WaitingForPlayers State = iota
	CardsDealt
	Finished
)

type Outcome int

const (
	Undecided Outcome = iota
	Win
	Lose
	Push
)

type Action int

const (
	Hit Action = iota
	Stand
)

type Player struct {
	Id      string  `json:"-"`
	Name    string  `json:"name"`
	IsReady bool    `json:"isReady"`
	Chips   int     `json:"chips"`
	Bet     int     `json:"bet"`
	Outcome Outcome `json:"outcome"`
}

type Blackjack struct {
	Deck          []deck.Card   `json:"-"`
	Hands         [][]deck.Card `json:"hands"`
	Players       []*Player     `json:"players"`
	State         State         `json:"state"`
	CurrentPlayer int           `json:"currentPlayer"`
}

func NewPlayer(id string, name string) Player {
	return Player{
		Id:      id,
		Name:    name,
		IsReady: false,
		Chips:   initialChips,
		Bet:     0,
		Outcome: Undecided,
	}
}

func New() Blackjack {
	dealerHand := []deck.Card{}
	return Blackjack{
		Deck:          deck.New(deck.WithShuffle()),
		Hands:         [][]deck.Card{dealerHand},
		Players:       []*Player{},
		State:         WaitingForPlayers,
		CurrentPlayer: 1,
	}
}

func (b *Blackjack) AddPlayer(id string, name string) {
	if b.State != WaitingForPlayers {
		return
	}
	newPlayer := NewPlayer(id, name)
	b.Players = append(b.Players, &newPlayer)
	b.Hands = append(b.Hands, []deck.Card{})
}

func (b *Blackjack) TogglePlayerReady(playerId string) {
	if b.State != WaitingForPlayers {
		return
	}
	allPlayersReady := true
	for _, player := range b.Players {
		if player.Id == playerId {
			player.IsReady = !player.IsReady
		}
		if !player.IsReady {
			allPlayersReady = false
		}
	}
	if allPlayersReady {
		b.Deal()
		b.State = CardsDealt
	}
}

func (b *Blackjack) Deal() {
	if b.State == CardsDealt {
		return
	}
	for range 2 {
		for player := 0; player < len(b.Hands); player++ {
			b.Hands[player] = append(b.Hands[player], b.Deck[0])
			b.Deck = b.Deck[1:]
		}
	}
}

func (b *Blackjack) GetPlayerCount() int {
	return len(b.Hands) - 1
}

func (b *Blackjack) GetPlayerHand(playerNumber int) []deck.Card {
	return b.Hands[playerNumber+1]
}

func (b *Blackjack) GetDealerHand() []deck.Card {
	return b.Hands[0]
}

func (b *Blackjack) PlayerAction(playerNumber int, action Action) bool {
	if b.State != CardsDealt || playerNumber+1 != b.CurrentPlayer {
		return false
	}

	switch action {
	case Hit:
		b.Hands[playerNumber+1] = append(b.Hands[playerNumber+1], b.Deck[0])
		b.Deck = b.Deck[1:]
		b.CurrentPlayer++
	case Stand:
		b.CurrentPlayer++
	}

	// Dealer's turn
	if b.CurrentPlayer == len(b.Hands) { // TODO(fix): CurrentPlayer should start from 0
		// TODO(feat): Implement dealer AI
		b.State = Finished
		b.DetermineOutcomes()
	}

	return true
}

func (b *Blackjack) DetermineOutcomes() {
	if b.State != Finished {
		return
	}
	for i := range b.GetPlayerCount() {
		b.Players[i].Outcome = determineOutcome(b.GetDealerHand(), b.GetPlayerHand(i))
	}
}

// nolint: mnd
func determineOutcome(dealerHand []deck.Card, playerHand []deck.Card) Outcome {
	dealerScore := getScore(dealerHand)
	playerScore := getScore(playerHand)
	dealerHasBlackjack := isBlackjack(dealerHand)
	playerHasBlackjack := isBlackjack(playerHand)

	if playerHasBlackjack && dealerHasBlackjack {
		return Push
	}
	if playerHasBlackjack {
		return Win
	}
	if dealerHasBlackjack {
		return Lose
	}
	if playerScore > 21 {
		return Lose
	}
	if dealerScore > 21 {
		return Win
	}
	if playerScore > dealerScore {
		return Win
	}
	if playerScore < dealerScore {
		return Lose
	}
	return Push
}

func getScore(hand []deck.Card) int {
	score := 0
	aceCount := 0
	for _, card := range hand {
		if card.Rank == deck.Ace {
			aceCount++
		}
		score += int(card.Rank)
	}
	for aceCount > 0 {
		if score > 21 { //nolint: mnd
			score -= 10
			aceCount--
		} else {
			break
		}
	}
	return score
}

func isBlackjack(hand []deck.Card) bool {
	if len(hand) != 2 { //nolint: mnd
		return false
	}
	isAce := func(i int) bool {
		return hand[i].Rank == deck.Ace
	}
	isTenOrQKJ := func(i int) bool {
		return hand[i].Rank == deck.Ten ||
			hand[i].Rank == deck.Jack ||
			hand[i].Rank == deck.Queen ||
			hand[i].Rank == deck.King
	}
	return (isAce(0) && isTenOrQKJ(1)) || (isAce(1) && isTenOrQKJ(0))
}
