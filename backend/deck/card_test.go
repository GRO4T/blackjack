package deck_test

import (
	"testing"

	"github.com/GRO4T/bjackapi/deck"
)

func AssertDeckSorted(t *testing.T, d []deck.Card) {
	t.Helper()
	for i := 1; i < len(d); i++ {
		if d[i].Suit < d[i-1].Suit {
			t.Errorf("Deck is not sorted by Suit")
		}
		if d[i].Suit == d[i-1].Suit && d[i].Rank < d[i-1].Rank {
			t.Errorf("Deck is not sorted by Rank")
		}
	}
}

func TestNew(t *testing.T) {
	d := deck.New()
	if len(d) != 52 {
		t.Errorf("Expected deck length of 52, but got %v", len(d))
	}
	AssertDeckSorted(t, d)
}

func TestNewWithShuffle(t *testing.T) {
	d := deck.New(deck.WithShuffle())
	if len(d) != 52 {
		t.Errorf("Expected deck length of 52, but got %v", len(d))
	}
	sorted := true
	for i := 1; i < len(d); i++ {
		if d[i].Suit < d[i-1].Suit {
			sorted = false
			break
		}
		if d[i].Suit == d[i-1].Suit && d[i].Rank < d[i-1].Rank {
			sorted = false
			break
		}
	}
	if sorted {
		t.Errorf("Expected deck to be shuffled, but it was sorted")
	}
}

func TestNewWithJokers(t *testing.T) {
	d := deck.New(deck.WithJokers(1))
	if len(d) != 53 {
		t.Errorf("Expected deck length of 53, but got %v", len(d))
	}
	AssertDeckSorted(t, d)
}

func TestWithFilter(t *testing.T) {
	d := deck.New(deck.WithFilter([]deck.Card{{Suit: deck.Spades, Rank: deck.Ace}}))
	if len(d) != 51 {
		t.Errorf("Expected deck length of 51, but got %v", len(d))
	}
	for _, card := range d {
		if card.Suit == deck.Spades && card.Rank == deck.Ace {
			t.Errorf("Expected Ace of Spades to be filtered out")
		}
	}
}

func TestShuffle(t *testing.T) {
	// Arrange
	d := deck.New()
	oldDeck := append([]deck.Card{}, d...)
	// Act
	deck.Shuffle(d)
	// Assert
	equal := true
	for i := range d {
		if d[i] != oldDeck[i] {
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
	d := deck.New()
	deck.Shuffle(d)
	// Act
	deck.Sort(d)
	// Assert
	for i := 1; i < len(d); i++ {
		if d[i].Suit < d[i-1].Suit {
			t.Errorf("Deck is not sorted by Suit")
		}
		if d[i].Suit == d[i-1].Suit && d[i].Rank < d[i-1].Rank {
			t.Errorf("Deck is not sorted by Rank")
		}
	}
}

func TestSortWithCustomLess(t *testing.T) {
	// Arrange
	d := deck.New()
	deck.Shuffle(d)
	// Act
	deck.Sort(d, func(i, j int) bool {
		if d[i].Suit == d[j].Suit {
			return d[i].Rank > d[j].Rank
		}
		return d[i].Suit > d[j].Suit
	})
	// Assert
	for i := 1; i < len(d); i++ {
		if d[i].Suit > d[i-1].Suit {
			t.Errorf("Deck is not sorted by Suit")
		}
		if d[i].Suit == d[i-1].Suit && d[i].Rank > d[i-1].Rank {
			t.Errorf("Deck is not sorted by Rank")
		}
	}
}
