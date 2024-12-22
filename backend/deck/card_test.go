package deck

import (
	"testing"
)

func AssertDeckSorted(t *testing.T, deck []Card) {
	for i := 1; i < len(deck); i++ {
		if deck[i].Suit < deck[i-1].Suit {
			t.Errorf("Deck is not sorted by Suit")
		}
		if deck[i].Suit == deck[i-1].Suit && deck[i].Rank < deck[i-1].Rank {
			t.Errorf("Deck is not sorted by Rank")
		}
	}
}

func TestNew(t *testing.T) {
	deck := New()
	if len(deck) != 52 {
		t.Errorf("Expected deck length of 52, but got %v", len(deck))
	}
	AssertDeckSorted(t, deck)
}

func TestNewWithShuffle(t *testing.T) {
	deck := New(WithShuffle())
	if len(deck) != 52 {
		t.Errorf("Expected deck length of 52, but got %v", len(deck))
	}
	sorted := true
	for i := 1; i < len(deck); i++ {
		if deck[i].Suit < deck[i-1].Suit {
			sorted = false
			break
		}
		if deck[i].Suit == deck[i-1].Suit && deck[i].Rank < deck[i-1].Rank {
			sorted = false
			break
		}
	}
	if sorted {
		t.Errorf("Expected deck to be shuffled, but it was sorted")
	}
}

func TestNewWithJokers(t *testing.T) {
	deck := New(WithJokers(1))
	if len(deck) != 53 {
		t.Errorf("Expected deck length of 53, but got %v", len(deck))
	}
	AssertDeckSorted(t, deck)
}

func TestWithFilter(t *testing.T) {
	deck := New(WithFilter([]Card{Card{Suit: Spades, Rank: Ace}}))
	if len(deck) != 51 {
		t.Errorf("Expected deck length of 51, but got %v", len(deck))
	}
	for _, card := range deck {
		if card.Suit == Spades && card.Rank == Ace {
			t.Errorf("Expected Ace of Spades to be filtered out")
		}
	}
}

func TestShuffle(t *testing.T) {
	// Arrange
	deck := New()
	oldDeck := append([]Card{}, deck...)
	// Act
	Shuffle(deck)
	// Assert
	equal := true
	for i := range deck {
		if deck[i] != oldDeck[i] {
			equal = false
			break
		}
	}
	if equal {
		t.Errorf("Expected deck to be shuffled, but it was the same")
	}
}

func TestSort(t *testing.T) {
	// Arrange
	deck := New()
	Shuffle(deck)
	// Act
	Sort(deck)
	// Assert
	for i := 1; i < len(deck); i++ {
		if deck[i].Suit < deck[i-1].Suit {
			t.Errorf("Deck is not sorted by Suit")
		}
		if deck[i].Suit == deck[i-1].Suit && deck[i].Rank < deck[i-1].Rank {
			t.Errorf("Deck is not sorted by Rank")
		}
	}
}

func TestSortWithCustomLess(t *testing.T) {
	// Arrange
	deck := New()
	Shuffle(deck)
	// Act
	Sort(deck, func(i, j int) bool {
		if deck[i].Suit == deck[j].Suit {
			return deck[i].Rank > deck[j].Rank
		}
		return deck[i].Suit > deck[j].Suit
	})
	// Assert
	for i := 1; i < len(deck); i++ {
		if deck[i].Suit > deck[i-1].Suit {
			t.Errorf("Deck is not sorted by Suit")
		}
		if deck[i].Suit == deck[i-1].Suit && deck[i].Rank > deck[i-1].Rank {
			t.Errorf("Deck is not sorted by Rank")
		}
	}
}
