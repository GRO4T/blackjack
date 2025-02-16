// TODO: Document this package

package deck

import (
	"crypto/rand"
	"math/big"
	"sort"
)

//go:generate stringer -type=Rank
type Rank int

const (
	AllRanks Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Joker
)

//go:generate stringer -type=Suit
type Suit int

const (
	AllSuits Suit = iota
	Spades
	Diamonds
	Clubs
	Hearts
)

type Card struct {
	Rank Rank `json:"rank"`
	Suit Suit `json:"suit"`
}

func New(options ...func([]Card) []Card) []Card {
	deck := []Card{}
	for suit := Spades; suit <= Hearts; suit++ {
		for rank := Ace; rank <= King; rank++ {
			deck = append(deck, Card{Rank: rank, Suit: suit})
		}
	}
	for _, o := range options {
		deck = o(deck)
	}
	return deck
}

func Sort(deck []Card, less ...func(i, j int) bool) {
	if len(less) > 0 {
		sort.Slice(deck, less[0])
	} else {
		sort.Slice(deck, func(i, j int) bool {
			if deck[i].Suit == deck[j].Suit {
				return deck[i].Rank < deck[j].Rank
			}
			return deck[i].Suit < deck[j].Suit
		})
	}
}

func Shuffle(deck []Card) {
	// Fisher-Yates shuffle
	for i := len(deck) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			panic(err)
		}
		deck[i], deck[j.Int64()] = deck[j.Int64()], deck[i]
	}
}

func WithShuffle() func([]Card) []Card {
	return func(deck []Card) []Card {
		Shuffle(deck)
		return deck
	}
}

func WithJokers(n int) func([]Card) []Card {
	return func(deck []Card) []Card {
		for i := range n {
			deck = append(deck, Card{Rank: Joker, Suit: Suit(i)})
		}
		Sort(deck)
		return deck
	}
}

func match(card Card, pattern Card) bool {
	if pattern.Rank != AllRanks && card.Rank != pattern.Rank {
		return false
	}
	if pattern.Suit != AllSuits && card.Suit != pattern.Suit {
		return false
	}
	return true
}

func WithFilter(discardedCards []Card) func([]Card) []Card {
	return func(deck []Card) []Card {
		filtered := []Card{}
		for _, card := range deck {
			discard := false
			for _, discardedCard := range discardedCards {
				if match(card, discardedCard) {
					discard = true
					continue
				}
			}
			if !discard {
				filtered = append(filtered, card)
			}
		}
		return filtered
	}
}

func WithMultipleDecks(n int) func([]Card) []Card {
	return func(deck []Card) []Card {
		singleDeck := append([]Card{}, deck...)
		for i := 1; i < n; i++ {
			deck = append(deck, singleDeck...)
		}
		Sort(deck)
		return deck
	}
}
