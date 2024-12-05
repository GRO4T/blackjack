package main

import (
	"fmt"

	"github.com/GRO4T/blackjack/deck"
)

type Blackjack struct {
	Deck       []deck.Card
	PlayerHand []deck.Card
	DealerHand []deck.Card
}

func NewBlackjack() Blackjack {
	return Blackjack{
		Deck:       deck.New(deck.WithShuffle()),
		PlayerHand: []deck.Card{},
		DealerHand: []deck.Card{},
	}
}

func (b *Blackjack) Play() {
	b.deal()

	b.playerTurn()
	playerScore := getScore(b.PlayerHand)
	if playerScore > 21 {
		fmt.Println("Player busts! Dealer wins.")
		return
	}

	b.dealerTurn()
	dealerScore := getScore(b.DealerHand)
	if dealerScore > 21 {
		fmt.Println("Dealer busts! Player wins.")
		return
	}

	if isNaturalBlackjack(b.DealerHand) {
		fmt.Println("Dealer has a natural blackjack! Dealer wins.")
		return
	} else if isNaturalBlackjack(b.PlayerHand) {
		fmt.Println("Player has a natural blackjack! Player wins.")
		return
	}

	if playerScore > dealerScore {
		fmt.Println("Player wins!")
	} else if playerScore == dealerScore {
		fmt.Println("It's a tie!")
	} else {
		fmt.Println("Dealer wins!")
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

func isNaturalBlackjack(hand []deck.Card) bool {
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

func (b *Blackjack) deal() {
	b.PlayerHand = append(b.PlayerHand, b.Deck[0])
	b.DealerHand = append(b.DealerHand, b.Deck[1])
	b.PlayerHand = append(b.PlayerHand, b.Deck[2])
	b.DealerHand = append(b.DealerHand, b.Deck[3])
	b.Deck = b.Deck[4:]
}

func (b *Blackjack) playerTurn() {
	fmt.Printf("Player: %s (score=%d)\n", b.PlayerHand, getScore(b.PlayerHand))
	fmt.Println("1. Hit")
	fmt.Println("2. Stand")
	endTurn := false
	for !endTurn {
		var action int
		fmt.Scan(&action)
		if action == 1 {
			b.PlayerHand = append(b.PlayerHand, b.Deck[0])
			b.Deck = b.Deck[1:]
			fmt.Printf("Player: %s (score=%d)\n", b.PlayerHand, getScore(b.PlayerHand))
		} else if action != 2 {
			fmt.Println("Invalid action. Please choose again.")
			fmt.Println("1. Hit")
			fmt.Println("2. Stand")
			continue
		}
		endTurn = true
	}
}

func (b *Blackjack) dealerTurn() {
	fmt.Printf("Dealer: Here's my hand (score=%d):\n", getScore(b.DealerHand))
	for _, card := range b.DealerHand {
		fmt.Println(card)
	}
}
