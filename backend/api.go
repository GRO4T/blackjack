// TODO: Improve the file structure (types and functions are mixed up)

package main

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"

	"log/slog"
)

const maxId int64 = 1000
const maxPlayers = 1

type Api struct {
	Games map[string]*Blackjack
}

type CreateGameResponse struct {
	TableId string `json:"tableId"`
}

type AddPlayerResponse struct {
	PlayerId string `json:"playerId"`
}

func NewApi() Api {
	return Api{
		Games: map[string]*Blackjack{},
	}
}

func (a *Api) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	slog.Info("POST /tables\n") // TODO: Can we simply make this automatic?

	tableId := getRandomId()
	newGame := NewBlackjack()
	a.Games[tableId] = &newGame

	var resp CreateGameResponse
	resp.TableId = tableId
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Debug("Created a new game %s\n", tableId)
}

func (a *Api) GetGameState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	slog.Info("GET /tables/%s\n", r.PathValue("tableId"))

	tableId := r.PathValue("tableId")
	game, ok := a.Games[tableId]
	if !ok {
		slog.Error("Game not found: %s\n", tableId)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(game); err != nil {
		slog.Error("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Debug("Retrieved game state for %s\n", tableId)
}

func (a *Api) AddPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	slog.Info("POST /tables/players/%s\n", r.PathValue("tableId"))

	tableId := r.PathValue("tableId")
	game, ok := a.Games[tableId]
	if !ok {
		slog.Error("Game not found: %s\n", tableId)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	if len(game.Players) >= maxPlayers {
		slog.Error("Game is full: %s\n", tableId)
		http.Error(w, "Game is full", http.StatusForbidden)
		return
	}

	if game.State != WaitingForPlayers {
		slog.Error("Game has already started: %s\n", tableId)
		http.Error(w, "Game has already started", http.StatusForbidden)
		return
	}

	playerId := getRandomId()
	game.AddPlayer(playerId, "Bob")

	var resp AddPlayerResponse
	resp.PlayerId = playerId
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Debug("Added player %s to game %s\n", playerId, tableId)
}

func (a *Api) TogglePlayerReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	slog.Info("POST /tables/ready/%s/%s\n", r.PathValue("tableId"), r.PathValue("playerId"))

	tableId := r.PathValue("tableId")
	game, ok := a.Games[tableId]
	if !ok {
		slog.Error("Game not found: %s\n", tableId)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	playerId := r.PathValue("playerId")
	playerIndex := -1
	for i, p := range game.Players {
		if p.Id == playerId {
			playerIndex = i
			break
		}
	}
	if playerIndex == -1 {
		slog.Error("Player not found: %s\n", playerId)
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	if game.State != WaitingForPlayers {
		slog.Error("Game is not waiting for readiness: %s\n", tableId)
		http.Error(w, "Game is not waiting for readiness", http.StatusForbidden)
		return
	}

	game.TogglePlayerReady(playerId)
	player := game.Players[playerIndex]
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(player); err != nil {
		slog.Error("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Debug("Toggled readiness for player %s\n", playerId)
}

// nolint: cyclop
func (a *Api) PlayerAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	slog.Info("POST /tables/%s/%s?action=%s\n", r.PathValue("tableId"), r.PathValue("playerId"), r.URL.Query().Get("action"))

	tableId := r.PathValue("tableId")
	game, ok := a.Games[tableId]
	if !ok {
		slog.Error("Game not found: %s\n", tableId)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	playerId := r.PathValue("playerId")
	playerIndex := -1
	for i, p := range game.Players {
		if p.Id == playerId {
			playerIndex = i
			break
		}
	}
	if playerIndex == -1 {
		slog.Error("Player not found: %s\n", playerId)
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	action := r.URL.Query().Get("action")

	switch action {
	case "hit":
		if !game.PlayerAction(playerIndex, Hit) {
			http.Error(w, "Invalid action", http.StatusBadRequest)
			return
		}
		slog.Debug("Player %s hit\n", playerId)
	case "stand":
		if !game.PlayerAction(playerIndex, Stand) {
			http.Error(w, "Invalid action", http.StatusBadRequest)
			return
		}
		slog.Debug("Player %s stood\n", playerId)
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}
}

func getRandomId() string {
	id, err := rand.Int(rand.Reader, big.NewInt(maxId))
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(id.Int64()))
}
