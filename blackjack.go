// TODO: Change return types to (<type>, error) and handle errors
package main

import (
	"github.com/GRO4T/blackjack/deck"
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
	IsReady bool    `json:"readiness"`
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
		Chips:   100,
		Bet:     0,
		Outcome: Undecided,
	}
}

func NewBlackjack() Blackjack {
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
	for cardCount := 0; cardCount < 2; cardCount++ {
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
	if b.CurrentPlayer == len(b.Hands) {
		// TODO: Implement dealer AI
		b.State = Finished
		b.DetermineOutcomes()
	}

	return true
}

func (b *Blackjack) DetermineOutcomes() {
	if b.State != Finished {
		return
	}

	dealerScore := getScore(b.GetDealerHand())
	dealerHasBlackjack := isBlackjack(b.GetDealerHand())
	for i := 0; i < b.GetPlayerCount(); i++ {
		player := b.Players[i]
		playerScore := getScore(b.GetPlayerHand(i))
		playerHasBlackjack := isBlackjack(b.GetPlayerHand(i))

		if playerHasBlackjack && dealerHasBlackjack {
			player.Outcome = Push
		} else if playerHasBlackjack {
			player.Outcome = Win
		} else if dealerHasBlackjack {
			player.Outcome = Lose
		} else if playerScore > 21 {
			player.Outcome = Lose
		} else if dealerScore > 21 {
			player.Outcome = Win
		} else if playerScore > dealerScore {
			player.Outcome = Win
		} else if playerScore < dealerScore {
			player.Outcome = Lose
		} else {
			player.Outcome = Push
		}
	}
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
		if score > 21 {
			score -= 10
			aceCount--
		} else {
			break
		}
	}
	return score
}

func isBlackjack(hand []deck.Card) bool {
	if len(hand) != 2 {
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
