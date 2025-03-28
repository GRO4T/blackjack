// nolint: noctx
package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GRO4T/bjack-api/blackjack"
	"github.com/GRO4T/bjack-api/rest"
)

func TestCreateGame(t *testing.T) {
	// Arrange
	api := rest.NewApi()
	server := httptest.NewServer(http.HandlerFunc(api.CreateGame))
	body := rest.CreateGameRequest{PlayerName: "Player 1"}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	resp, err := http.Post(server.URL, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", resp.Status)
	}
	if len(api.Games) == 0 {
		t.Fatal("Game not created")
	}
	for _, game := range api.Games {
		if len(game.Players) > 0 {
			t.Errorf("Expected 0 players; got %v", len(game.Players))
		}
		if len(game.Hands) != 1 { // Game starts with a dealer hand
			t.Errorf("Expected 1 hands; got %v", len(game.Hands))
		}
	}
}

func TestGetGameState(t *testing.T) {
	// Arrange
	api := rest.NewApi()
	game := blackjack.New(nil)
	api.Games["1"] = &game

	// Act
	request, err := http.NewRequest(http.MethodGet, "/tables/{tableId}", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.SetPathValue("tableId", "1")
	responseWriter := httptest.NewRecorder()
	api.GetGameState(responseWriter, request)
	resp := responseWriter.Result()
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", resp.Status)
	}
	var gameState blackjack.Blackjack
	if err := json.NewDecoder(resp.Body).Decode(&gameState); err != nil {
		t.Fatal(err)
	}
}

func buildAddPlayerRequest(t *testing.T, tableId string, playerName string) *http.Request {
	t.Helper()
	body := rest.AddPlayerRequest{PlayerName: playerName}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	request, err := http.NewRequest(http.MethodPost, "/tables/players/{tableId}", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	request.SetPathValue("tableId", tableId)
	return request
}

func buildRemovePlayerRequest(t *testing.T, tableId string, playerId string) *http.Request {
	t.Helper()
	request, err := http.NewRequest(http.MethodDelete, "/tables/players/{tableId}/{playerId}", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.SetPathValue("tableId", tableId)
	request.SetPathValue("playerId", playerId)
	return request
}

func TestAddPlayer(t *testing.T) {
	// Arrange
	api := rest.NewApi()
	game := blackjack.New(nil)
	api.Games["1"] = &game
	responseWriter := httptest.NewRecorder()
	request := buildAddPlayerRequest(t, "1", "Player 1")

	// Act
	api.AddPlayer(responseWriter, request)
	resp := responseWriter.Result()
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", resp.Status)
	}
	if len(api.Games["1"].Players) != 1 {
		t.Errorf("Expected 1 player; got %v", len(api.Games["1"].Players))
	}
}

func TestRemovePlayer(t *testing.T) {
	// Arrange
	api := rest.NewApi()
	game := blackjack.New(nil)
	api.Games["1"] = &game

	newPlayer, _ := game.AddPlayer("Player 1")

	removePlayerResponseWriter := httptest.NewRecorder()
	removePlayerRequest := buildRemovePlayerRequest(t, "1", newPlayer.Id)

	// Act
	api.RemovePlayer(removePlayerResponseWriter, removePlayerRequest)
	resp := removePlayerResponseWriter.Result()
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", resp.Status)
	}
	if len(api.Games["1"].Players) != 0 {
		t.Errorf("Expected 0 players; got %v", len(api.Games["1"].Players))
	}
}

func TestRemovePlayerWhenGameAlreadyStarted(t *testing.T) {
	// Arrange
	api := rest.NewApi()
	game := blackjack.New(nil)
	api.Games["1"] = &game

	newPlayer, _ := game.AddPlayer("Player 1")
	_, err := game.TogglePlayerReady(newPlayer.Id)
	if err != nil {
		t.Fatal(err)
	}

	removePlayerResponseWriter := httptest.NewRecorder()
	removePlayerRequest := buildRemovePlayerRequest(t, "1", newPlayer.Id)

	// Act
	api.RemovePlayer(removePlayerResponseWriter, removePlayerRequest)
	resp := removePlayerResponseWriter.Result()
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status 400; got %v\n", resp.Status)
	}
	if len(game.Players) != 1 {
		t.Errorf("Expected 1 players; got %v", len(game.Players))
	}
}

func TestTogglePlayerReadyWhenPlayerNotReady(t *testing.T) {
	// Arrange
	api := rest.NewApi()
	game := blackjack.New(nil)
	newPlayer, _ := game.AddPlayer("Player 1")
	api.Games["1"] = &game

	// Act
	request, err := http.NewRequest(http.MethodPost, "/tables/ready/{tableId}/{playerId}", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.SetPathValue("tableId", "1")
	request.SetPathValue("playerId", newPlayer.Id)
	responseWriter := httptest.NewRecorder()
	api.TogglePlayerReady(responseWriter, request)
	resp := responseWriter.Result()
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", resp.Status)
	}
	var player blackjack.Player
	if err := json.NewDecoder(resp.Body).Decode(&player); err != nil {
		t.Fatal(err)
	}
	if !player.IsReady {
		t.Error("Expected player to be ready")
	}
	if game.State != blackjack.CardsDealt {
		t.Error("Expected game to be in CardsDealt state")
	}
}

func TestTogglePlayerReadyWhenPlayerReady(t *testing.T) {
	// Arrange
	api := rest.NewApi()
	game := blackjack.New(nil)
	newPlayer, _ := game.AddPlayer("Player 1")
	game.State = blackjack.WaitingForPlayers
	api.Games["1"] = &game
	api.Games["1"].Players[0].IsReady = true

	// Act
	request, err := http.NewRequest(http.MethodPost, "/tables/ready/{tableId}/{playerId}", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.SetPathValue("tableId", "1")
	request.SetPathValue("playerId", newPlayer.Id)
	responseWriter := httptest.NewRecorder()
	api.TogglePlayerReady(responseWriter, request)
	resp := responseWriter.Result()
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", resp.Status)
	}
	var player blackjack.Player
	if err := json.NewDecoder(resp.Body).Decode(&player); err != nil {
		t.Fatal(err)
	}
	if player.IsReady {
		t.Error("Expected player to be not ready")
	}
	if game.State != blackjack.WaitingForPlayers {
		t.Error("Expected game to be in WaitingForPlayers state")
	}
}

func TestPlayerHit(t *testing.T) {
	// Arrange
	api := rest.NewApi()
	game := blackjack.New(nil)
	newPlayer, _ := game.AddPlayer("Player 1")
	err := game.Deal()
	if err != nil {
		t.Fatal(err)
	}
	game.State = blackjack.CardsDealt
	api.Games["1"] = &game

	// Act
	request, err := http.NewRequest(http.MethodPost, "/tables/{tableId}/{playerId}?action=hit", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.SetPathValue("tableId", "1")
	request.SetPathValue("playerId", newPlayer.Id)
	responseWriter := httptest.NewRecorder()
	api.PlayerAction(responseWriter, request)
	resp := responseWriter.Result()
	defer resp.Body.Close()

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", resp.Status)
	}
	if len(api.Games["1"].GetPlayerHand(0)) != 3 {
		t.Errorf("Expected 3 cards; got %v", len(api.Games["1"].GetPlayerHand(0)))
	}
}

//nolint:cyclop
func TestSimpleGame(t *testing.T) {
	api := rest.NewApi()

	// Create game
	createGameBody := rest.CreateGameRequest{PlayerName: "Player 1"}
	createGameBodyBytes, err := json.Marshal(createGameBody)
	if err != nil {
		t.Fatal(err)
	}
	createGameRequest, err := http.NewRequest(http.MethodPost, "/tables", bytes.NewReader(createGameBodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	createGameResponseWriter := httptest.NewRecorder()
	api.CreateGame(createGameResponseWriter, createGameRequest)
	createGameResp := createGameResponseWriter.Result()
	defer createGameResp.Body.Close()
	if createGameResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", createGameResp.Status)
	}
	var createGameRespBody rest.CreateGameResponse
	if err := json.NewDecoder(createGameResp.Body).Decode(&createGameRespBody); err != nil {
		t.Fatal(err)
	}
	tableId := createGameRespBody.TableId

	// Add player
	addPlayerBody := rest.AddPlayerRequest{PlayerName: "Player 1"}
	addPlayerBodyBytes, err := json.Marshal(addPlayerBody)
	if err != nil {
		t.Fatal(err)
	}
	addPlayerRequest, err := http.NewRequest(http.MethodPost, "/tables/players/{tableId}", bytes.NewReader(addPlayerBodyBytes))
	addPlayerRequest.SetPathValue("tableId", tableId)
	if err != nil {
		t.Fatal(err)
	}
	addPlayerResponseWriter := httptest.NewRecorder()
	api.AddPlayer(addPlayerResponseWriter, addPlayerRequest)
	addPlayerResp := addPlayerResponseWriter.Result()
	defer addPlayerResp.Body.Close()
	if addPlayerResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", addPlayerResp.Status)
	}
	var addPlayerRespBody rest.AddPlayerResponse
	if err := json.NewDecoder(addPlayerResp.Body).Decode(&addPlayerRespBody); err != nil {
		t.Fatal(err)
	}
	playerId := addPlayerRespBody.PlayerId

	// Toggle player ready
	togglePlayerReadyRequest, err := http.NewRequest(http.MethodPost, "/tables/ready/{tableId}/{playerId}", nil)
	togglePlayerReadyRequest.SetPathValue("tableId", tableId)
	togglePlayerReadyRequest.SetPathValue("playerId", playerId)
	if err != nil {
		t.Fatal(err)
	}
	togglePlayerReadyResponseWriter := httptest.NewRecorder()
	api.TogglePlayerReady(togglePlayerReadyResponseWriter, togglePlayerReadyRequest)
	togglePlayerReadyResp := togglePlayerReadyResponseWriter.Result()
	defer togglePlayerReadyResp.Body.Close()
	if togglePlayerReadyResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", togglePlayerReadyResp.Status)
	}

	// Player hit
	playerHitRequest, err := http.NewRequest(http.MethodPost, "/tables/{tableId}/{playerId}?action=hit", nil)
	playerHitRequest.SetPathValue("tableId", tableId)
	playerHitRequest.SetPathValue("playerId", playerId)
	if err != nil {
		t.Fatal(err)
	}
	playerHitResponseWriter := httptest.NewRecorder()
	api.PlayerAction(playerHitResponseWriter, playerHitRequest)
	playerHitResp := playerHitResponseWriter.Result()
	defer playerHitResp.Body.Close()
	if playerHitResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", playerHitResp.Status)
	}

	// Check game outcome
	getGameStateRequest, err := http.NewRequest(http.MethodGet, "/tables/{tableId}", nil)
	getGameStateRequest.SetPathValue("tableId", tableId)
	if err != nil {
		t.Fatal(err)
	}
	getGameStateResponseWriter := httptest.NewRecorder()
	api.GetGameState(getGameStateResponseWriter, getGameStateRequest)
	getGameStateResp := getGameStateResponseWriter.Result()
	defer getGameStateResp.Body.Close()
	if getGameStateResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v\n", getGameStateResp.Status)
	}
	var gameState blackjack.Blackjack
	if err := json.NewDecoder(getGameStateResp.Body).Decode(&gameState); err != nil {
		t.Fatal(err)
	}
	if gameState.Players[0].Outcome == blackjack.Undecided {
		t.Error("Expected player outcome to be decided")
	}
}
