package blackjack_test

import (
	"errors"
	"testing"

	"github.com/GRO4T/bjack-api/blackjack"
	"github.com/GRO4T/bjack-api/constant"
	"github.com/GRO4T/bjack-api/deck"
)

func TestNewPlayer(t *testing.T) {
	player := blackjack.NewPlayer("id1", "Alice")
	if player.Id != "id1" {
		t.Fatalf("Expected player ID 'id1'; got %v\n", player.Id)
	}
	if player.Name != "Alice" {
		t.Fatalf("Expected player Name 'Alice'; got %v\n", player.Name)
	}
	if player.IsReady {
		t.Fatalf("Expected IsReady false; got %v\n", player.IsReady)
	}
	if player.Chips != 100 {
		t.Fatalf("Expected 100 chips; got %v\n", player.Chips)
	}
	if player.Bet != 0 {
		t.Fatalf("Expected bet 0; got %v\n", player.Bet)
	}
	if player.Outcome != blackjack.Undecided {
		t.Fatalf("Expected outcome Undecided; got %v\n", player.Outcome)
	}
}

func TestCreateGame(t *testing.T) {
	b := blackjack.New(nil)
	if b.State != blackjack.WaitingForPlayers {
		t.Fatalf("Expected state %v; got %v\n", blackjack.WaitingForPlayers, b.State)
	}
}

func TestAddPlayer(t *testing.T) {
	b := blackjack.New(nil)
	_, err := b.AddPlayer("Player")
	if err != nil {
		t.Fatalf("Failed to add player: %v\n", err)
	}
	if len(b.Players) != 1 {
		t.Fatalf("Expected 1 player; got %v\n", len(b.Players))
	}
	if len(b.Hands) != 2 {
		t.Fatalf("Expected 2 hands (dealer+player); got %v\n", len(b.Hands))
	}
}

func TestAddPlayerWhenNameAlreadyTakenShouldFail(t *testing.T) {
	b := blackjack.New(nil)
	_, err := b.AddPlayer("Player")
	if err != nil {
		t.Fatalf("Failed to add player: %v\n", err)
	}
	_, err = b.AddPlayer("Player")
	if err == nil {
		t.Fatal("Game allowed to register multiple players with the same name\n")
	}
}

func TestAddPlayerWhenGameIsFullShouldFail(t *testing.T) {
	b := blackjack.New(nil)
	_, err := b.AddPlayer("Player")
	if err != nil {
		t.Fatalf("Failed to add player: %v\n", err)
	}
	_, err = b.AddPlayer("Player2")
	if err != nil {
		t.Fatalf("Failed to add player: %v\n", err)
	}
	_, err = b.AddPlayer("Player3")
	if err != nil {
		t.Fatalf("Failed to add player: %v\n", err)
	}
	_, err = b.AddPlayer("Player4")
	if !errors.Is(err, blackjack.ErrGameIsFull) {
		t.Fatalf("Expected ErrGameIsFull, got %v (constant.MaxPlayers is %v)\n", err, constant.MaxPlayers)
	}
}

func TestAddPlayerAfterGameStartedShouldFail(t *testing.T) {
	b := blackjack.New(nil)
	player, _ := b.AddPlayer("Player")
	_, _ = b.TogglePlayerReady(player.Id)
	_, err := b.AddPlayer("Player2")
	if !errors.Is(err, blackjack.ErrGameAlreadyStarted) {
		t.Fatalf("Expected ErrGameAlreadyStarted, got %v\n", err)
	}
}

func TestRemovePlayer(t *testing.T) {
	b := blackjack.New(nil)
	player, _ := b.AddPlayer("Player")
	err := b.RemovePlayer(player.Id)
	if err != nil {
		t.Fatalf("Failed to remove player: %v\n", err)
	}
	if len(b.Players) != 0 {
		t.Fatalf("Expected 0 players; got %v\n", len(b.Players))
	}
}

func TestRemovePlayerThatNotExistShouldFail(t *testing.T) {
	b := blackjack.New(nil)
	err := b.RemovePlayer("non-existent-id")
	if !errors.Is(err, blackjack.ErrNotFound) {
		t.Fatalf("Expected ErrNotFound, got %v\n", err)
	}
}

func TestRemovePlayerAfterGameStartedShouldFail(t *testing.T) {
	b := blackjack.New(nil)
	player, _ := b.AddPlayer("Player")
	_, err := b.TogglePlayerReady(player.Id)
	if err != nil {
		t.Fatalf("Failed to toggle player ready: %v\n", err)
	}
	err = b.RemovePlayer(player.Id)
	if !errors.Is(err, blackjack.ErrGameAlreadyStarted) {
		t.Fatalf("Expected ErrGameAlreadyStarted, got %v\n", err)
	}
}

func TestTogglePlayerReady(t *testing.T) {
	b := blackjack.New(nil)
	p1, _ := b.AddPlayer("Player1")
	p2, _ := b.AddPlayer("Player2")

	resP1, err := b.TogglePlayerReady(p1.Id)
	if err != nil {
		t.Fatalf("Failed to toggle player ready: %v\n", err)
	}
	if !resP1.IsReady {
		t.Fatalf("Expected player 1 to be ready\n")
	}
	if b.State != blackjack.WaitingForPlayers {
		t.Fatalf("Expected state WaitingForPlayers, got %v\n", b.State)
	}

	resP2, err := b.TogglePlayerReady(p2.Id)
	if err != nil {
		t.Fatalf("Failed to toggle player ready: %v\n", err)
	}
	if !resP2.IsReady {
		t.Fatalf("Expected player 2 to be ready\n")
	}
	if b.State != blackjack.CardsDealt {
		t.Fatalf("Expected state CardsDealt after all ready, got %v\n", b.State)
	}
	if len(b.GetDealerHand()) != 2 {
		t.Fatalf("Expected 2 cards in dealer hand; got %v\n", len(b.GetDealerHand()))
	}
	if len(b.GetPlayerHand(0)) != 2 || len(b.GetPlayerHand(1)) != 2 {
		t.Fatalf("Expected 2 cards in player hands\n")
	}
}

func TestTogglePlayerReadyWhenPlayerNotFoundShouldFail(t *testing.T) {
	b := blackjack.New(nil)
	_, err := b.TogglePlayerReady("invalid-id")
	if !errors.Is(err, blackjack.ErrNotFound) {
		t.Fatalf("Expected ErrNotFound, got %v\n", err)
	}
}

func TestTogglePlayerReadyAfterGameStartedShouldFail(t *testing.T) {
	b := blackjack.New(nil)
	p1, _ := b.AddPlayer("Player1")
	_, _ = b.TogglePlayerReady(p1.Id)

	_, err := b.TogglePlayerReady(p1.Id)
	if !errors.Is(err, blackjack.ErrGameAlreadyStarted) {
		t.Fatalf("Expected ErrGameAlreadyStarted, got %v\n", err)
	}
}

func TestDealWhenStateCardsDealtShouldFail(t *testing.T) {
	b := blackjack.New(nil)
	p1, _ := b.AddPlayer("Player1")
	_, _ = b.TogglePlayerReady(p1.Id)

	err := b.Deal()
	if !errors.Is(err, blackjack.ErrCardsAlreadyDealt) {
		t.Fatalf("Expected ErrCardsAlreadyDealt, got %v\n", err)
	}
}

func TestDealSuccess(t *testing.T) {
	b := blackjack.New(nil)
	_, _ = b.AddPlayer("Player1")
	err := b.Deal()
	if err != nil {
		t.Fatalf("Unexpected error on Deal: %v\n", err)
	}
	if len(b.GetDealerHand()) != 2 {
		t.Fatalf("Expected dealer hand to have 2 cards, got %v\n", len(b.GetDealerHand()))
	}
	if len(b.GetPlayerHand(0)) != 2 {
		t.Fatalf("Expected player hand to have 2 cards, got %v\n", len(b.GetPlayerHand(0)))
	}
}

func TestGetPlayerCountAndHands(t *testing.T) {
	b := blackjack.New(nil)
	_, _ = b.AddPlayer("Player1")
	_, _ = b.AddPlayer("Player2")

	if count := b.GetPlayerCount(); count != 2 {
		t.Fatalf("Expected GetPlayerCount() == 2; got %v\n", count)
	}

	_ = b.Deal()
	if len(b.GetDealerHand()) != 2 {
		t.Fatalf("Expected dealer hand len 2; got %v\n", len(b.GetDealerHand()))
	}
	if len(b.GetPlayerHand(0)) != 2 {
		t.Fatalf("Expected player 0 hand len 2; got %v\n", len(b.GetPlayerHand(0)))
	}
	if len(b.GetPlayerHand(1)) != 2 {
		t.Fatalf("Expected player 1 hand len 2; got %v\n", len(b.GetPlayerHand(1)))
	}
}

func TestPlayerActionWhenGameNotInProgressShouldFail(t *testing.T) {
	b := blackjack.New(nil)
	p1, _ := b.AddPlayer("Player1")
	err := b.PlayerAction(p1.Id, blackjack.Hit)
	if !errors.Is(err, blackjack.ErrGameNotInProgress) {
		t.Fatalf("Expected ErrGameNotInProgress, got %v\n", err)
	}
}

func TestPlayerActionWhenPlayerNotFoundShouldFail(t *testing.T) {
	b := blackjack.New(nil)
	p1, _ := b.AddPlayer("Player1")
	_, _ = b.TogglePlayerReady(p1.Id)

	err := b.PlayerAction("invalid-id", blackjack.Hit)
	if !errors.Is(err, blackjack.ErrNotFound) {
		t.Fatalf("Expected ErrNotFound, got %v\n", err)
	}
}

func TestPlayerActionWhenNotPlayerTurnShouldFail(t *testing.T) {
	b := blackjack.New(nil)
	_, _ = b.AddPlayer("Player1")
	p2, _ := b.AddPlayer("Player2")
	b.Players[0].IsReady = true
	_, _ = b.TogglePlayerReady(p2.Id)

	err := b.PlayerAction(p2.Id, blackjack.Hit)
	if !errors.Is(err, blackjack.ErrOtherPlayerTurn) {
		t.Fatalf("Expected ErrOtherPlayerTurn, got %v\n", err)
	}
}

func TestPlayerActionHitAndStand(t *testing.T) {
	b := blackjack.New(nil)
	p1, _ := b.AddPlayer("Player1")
	p2, _ := b.AddPlayer("Player2")
	b.Players[0].IsReady = true
	_, _ = b.TogglePlayerReady(p2.Id)

	initialDeckLen := len(b.Deck)

	err := b.PlayerAction(p1.Id, blackjack.Hit)
	if err != nil {
		t.Fatalf("Failed Hit action: %v\n", err)
	}
	if len(b.GetPlayerHand(0)) != 3 {
		t.Fatalf("Expected player 1 hand to have 3 cards after Hit; got %v\n", len(b.GetPlayerHand(0)))
	}
	if len(b.Deck) != initialDeckLen-1 {
		t.Fatalf("Expected deck to decrease by 1 card after Hit\n")
	}

	err = b.PlayerAction(p2.Id, blackjack.Stand)
	if err != nil {
		t.Fatalf("Failed Stand action: %v\n", err)
	}
	if len(b.GetPlayerHand(1)) != 2 {
		t.Fatalf("Expected player 2 hand to remain 2 cards after Stand; got %v\n", len(b.GetPlayerHand(1)))
	}

	if b.State != blackjack.Finished {
		t.Fatalf("Expected game State to be Finished; got %v\n", b.State)
	}
}

func TestOnStateChangedCallback(t *testing.T) {
	callCount := 0
	callback := func() {
		callCount++
	}

	b := blackjack.New(callback)

	p1, err := b.AddPlayer("Player1")
	if err != nil {
		t.Fatalf("AddPlayer failed: %v\n", err)
	}
	if callCount != 1 {
		t.Fatalf("Expected callCount 1 after AddPlayer, got %d\n", callCount)
	}

	_, err = b.TogglePlayerReady(p1.Id)
	if err != nil {
		t.Fatalf("TogglePlayerReady failed: %v\n", err)
	}
	if callCount != 2 {
		t.Fatalf("Expected callCount 2 after TogglePlayerReady, got %d\n", callCount)
	}

	err = b.PlayerAction(p1.Id, blackjack.Stand)
	if err != nil {
		t.Fatalf("PlayerAction failed: %v\n", err)
	}
	if callCount != 3 {
		t.Fatalf("Expected callCount 3 after PlayerAction, got %d\n", callCount)
	}

	b.State = blackjack.WaitingForPlayers
	err = b.RemovePlayer(p1.Id)
	if err != nil {
		t.Fatalf("RemovePlayer failed: %v\n", err)
	}
	if callCount != 4 {
		t.Fatalf("Expected callCount 4 after RemovePlayer, got %d\n", callCount)
	}
}

func TestDetermineOutcomesWhenNotFinished(t *testing.T) {
	b := blackjack.New(nil)
	p1, _ := b.AddPlayer("Player1")
	b.DetermineOutcomes()
	if p1.Outcome != blackjack.Undecided {
		t.Fatalf("Expected outcome Undecided when state not Finished; got %v\n", p1.Outcome)
	}
}

func TestDetermineOutcomesScenarios(t *testing.T) {
	tests := []struct {
		name           string
		dealerHand     []deck.Card
		playerHand     []deck.Card
		expectedResult blackjack.Outcome
	}{
		{
			name: "Both Blackjack -> Push",
			dealerHand: []deck.Card{
				{Rank: deck.Ace, Suit: deck.Spades},
				{Rank: deck.King, Suit: deck.Hearts},
			},
			playerHand: []deck.Card{
				{Rank: deck.Ace, Suit: deck.Clubs},
				{Rank: deck.Queen, Suit: deck.Diamonds},
			},
			expectedResult: blackjack.Push,
		},
		{
			name: "Player Blackjack vs Dealer 20 -> Win",
			dealerHand: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Spades},
				{Rank: deck.Ten, Suit: deck.Hearts},
			},
			playerHand: []deck.Card{
				{Rank: deck.Ace, Suit: deck.Clubs},
				{Rank: deck.Jack, Suit: deck.Diamonds},
			},
			expectedResult: blackjack.Win,
		},
		{
			name: "Dealer Blackjack vs Player 20 -> Lose",
			dealerHand: []deck.Card{
				{Rank: deck.Ace, Suit: deck.Spades},
				{Rank: deck.Jack, Suit: deck.Hearts},
			},
			playerHand: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Clubs},
				{Rank: deck.Ten, Suit: deck.Diamonds},
			},
			expectedResult: blackjack.Lose,
		},
		{
			name: "Player Bust (>21) -> Lose",
			dealerHand: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Spades},
				{Rank: deck.Seven, Suit: deck.Hearts},
			},
			playerHand: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Clubs},
				{Rank: deck.Eight, Suit: deck.Diamonds},
				{Rank: deck.Five, Suit: deck.Spades},
			},
			expectedResult: blackjack.Lose,
		},
		{
			name: "Dealer Bust (>21) -> Win",
			dealerHand: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Spades},
				{Rank: deck.Six, Suit: deck.Hearts},
				{Rank: deck.Eight, Suit: deck.Clubs},
			},
			playerHand: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Clubs},
				{Rank: deck.Seven, Suit: deck.Diamonds},
			},
			expectedResult: blackjack.Win,
		},
		{
			name: "Player higher score (20 vs 18) -> Win",
			dealerHand: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Spades},
				{Rank: deck.Eight, Suit: deck.Hearts},
			},
			playerHand: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Clubs},
				{Rank: deck.Ten, Suit: deck.Diamonds},
			},
			expectedResult: blackjack.Win,
		},
		{
			name: "Dealer higher score (19 vs 17) -> Lose",
			dealerHand: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Spades},
				{Rank: deck.Nine, Suit: deck.Hearts},
			},
			playerHand: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Clubs},
				{Rank: deck.Seven, Suit: deck.Diamonds},
			},
			expectedResult: blackjack.Lose,
		},
		{
			name: "Equal score non-blackjack (18 vs 18) -> Push",
			dealerHand: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Spades},
				{Rank: deck.Eight, Suit: deck.Hearts},
			},
			playerHand: []deck.Card{
				{Rank: deck.Nine, Suit: deck.Clubs},
				{Rank: deck.Nine, Suit: deck.Diamonds},
			},
			expectedResult: blackjack.Push,
		},
		{
			name: "Ace evaluated as 11 (Soft 20)",
			dealerHand: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Spades},
				{Rank: deck.Seven, Suit: deck.Hearts},
			},
			playerHand: []deck.Card{
				{Rank: deck.Ace, Suit: deck.Clubs},
				{Rank: deck.Nine, Suit: deck.Diamonds},
			},
			expectedResult: blackjack.Win,
		},
		{
			name: "Multiple Aces (Ace + Ace + Nine = 21)",
			dealerHand: []deck.Card{
				{Rank: deck.Ten, Suit: deck.Spades},
				{Rank: deck.Nine, Suit: deck.Hearts},
			},
			playerHand: []deck.Card{
				{Rank: deck.Ace, Suit: deck.Clubs},
				{Rank: deck.Ace, Suit: deck.Diamonds},
				{Rank: deck.Nine, Suit: deck.Spades},
			},
			expectedResult: blackjack.Win,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := blackjack.New(nil)
			p, err := b.AddPlayer("TestPlayer")
			if err != nil {
				t.Fatalf("AddPlayer failed: %v", err)
			}
			b.Hands = [][]deck.Card{tt.dealerHand, tt.playerHand}
			b.State = blackjack.Finished
			b.DetermineOutcomes()

			if p.Outcome != tt.expectedResult {
				t.Fatalf("Expected outcome %v; got %v", tt.expectedResult, p.Outcome)
			}
		})
	}
}
