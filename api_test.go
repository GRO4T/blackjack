package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateGame(t *testing.T) {
	// Arrange
	api := NewApi()
	server := httptest.NewServer(http.HandlerFunc(api.CreateGame))

	// Act
	resp, err := http.Post(server.URL, "application/json", nil)

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", resp.Status)
	}
	var gameState GameDto
	if err := json.NewDecoder(resp.Body).Decode(&gameState); err != nil {
		t.Fatal(err)
	}
	if len(api.Games) == 0 {
		t.Fatal("Game not created")
	}
	if len(gameState.Players) > 0 {
		t.Errorf("Expected 0 players; got %v", len(gameState.Players))
	}
	if len(gameState.Chips) > 0 {
		t.Errorf("Expected 0 chips; got %v", len(gameState.Chips))
	}
	if len(gameState.Bets) > 0 {
		t.Errorf("Expected 0 bets; got %v", len(gameState.Bets))
	}
	if len(gameState.Hands) != 0 {
		t.Errorf("Expected 0 hands; got %v", len(gameState.Hands))
	}
}

func TestGetGameState(t *testing.T) {
	// Arrange
	api := NewApi()
	api.Games["1"] = &Game{
		Table:   NewBlackjack(),
		TableId: "1",
		Players: []string{},
		Chips:   []int{},
		Bets:    []int{},
	}

	// Act
	request, err := http.NewRequest(http.MethodGet, "/tables/{tableId}", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.SetPathValue("tableId", "1")
	responseWriter := httptest.NewRecorder()
	api.GetGameState(responseWriter, request)
	resp := responseWriter.Result()

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", resp.Status)
	}
	var gameState GameDto
	if err := json.NewDecoder(resp.Body).Decode(&gameState); err != nil {
		t.Fatal(err)
	}
	if gameState.TableId != "1" {
		t.Errorf("Expected tableId 1; got %v", gameState.TableId)
	}
}

func TestAddPlayer(t *testing.T) {
	// Arrange
	api := NewApi()
	api.Games["1"] = &Game{
		Table:   NewBlackjack(),
		TableId: "1",
		Players: []string{},
		Chips:   []int{},
		Bets:    []int{},
	}

	// Act
	request, err := http.NewRequest(http.MethodPost, "/tables/players/{tableId}", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.SetPathValue("tableId", "1")
	responseWriter := httptest.NewRecorder()
	api.AddPlayer(responseWriter, request)
	resp := responseWriter.Result()

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", resp.Status)
	}
	if len(api.Games["1"].Players) != 1 {
		t.Errorf("Expected 1 player; got %v", len(api.Games["1"].Players))
	}
}

func TestPlayerHit(t *testing.T) {
	// Arrange
	api := NewApi()
	api.Games["1"] = &Game{
		Table:   NewBlackjack(),
		TableId: "1",
		Players: []string{"1"},
		Chips:   []int{100},
		Bets:    []int{0},
	}
	api.Games["1"].Table.deal()

	// Act
	request, err := http.NewRequest(http.MethodPost, "/tables/{tableId}/{playerId}?action=hit", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.SetPathValue("tableId", "1")
	request.SetPathValue("playerId", "1")
	responseWriter := httptest.NewRecorder()
	api.PlayerAction(responseWriter, request)
	resp := responseWriter.Result()

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", resp.Status)
	}
	if len(api.Games["1"].Table.PlayerHand) != 3 {
		t.Errorf("Expected 3 cards; got %v", len(api.Games["1"].Table.PlayerHand))
	}
}
