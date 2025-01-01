// TODO: Improve the file structure (types and functions are mixed up)

package rest

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/GRO4T/blackjack/blackjack"
	"github.com/GRO4T/blackjack/constant"
)

type RestApi struct {
	Games map[string]*blackjack.Blackjack
}

type CreateGameResponse struct {
	TableId string `json:"tableId"`
}

type AddPlayerResponse struct {
	PlayerId string `json:"playerId"`
}

func NewApi() RestApi {
	return RestApi{
		Games: map[string]*blackjack.Blackjack{},
	}
}

func (a *RestApi) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	tableId := getRandomId()
	newGame := blackjack.New()
	a.Games[tableId] = &newGame

	var resp CreateGameResponse
	resp.TableId = tableId
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(fmt.Sprintf("Failed to encode response: %v", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Debug("Created a new game", "tableId", tableId)
}

func (a *RestApi) GetGameState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	tableId := r.PathValue("tableId")

	game, ok := a.Games[tableId]
	if !ok {
		slog.Error("Game not found", "tableId", tableId)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(game); err != nil {
		slog.Error(fmt.Sprintf("Failed to encode response: %v", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Debug("Retrieved game state", "tableId", tableId)
}

func (a *RestApi) AddPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	tableId := r.PathValue("tableId")

	game, ok := a.Games[tableId]
	if !ok {
		slog.Error("Game not found", "tableId", tableId)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	// TODO(refactor): Move this check to the game object
	if len(game.Players) >= constant.MaxPlayers {
		slog.Error("Game is full", "tableId", tableId)
		http.Error(w, "Game is full", http.StatusForbidden)
		return
	}

	if game.State != blackjack.WaitingForPlayers {
		slog.Error("Game has already started", "tableId", tableId)
		http.Error(w, "Game has already started", http.StatusForbidden)
		return
	}

	playerId := getRandomId()
	game.AddPlayer(playerId, "Bob") // TODO: Allow to set player name in request

	var resp AddPlayerResponse
	resp.PlayerId = playerId
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(fmt.Sprintf("Failed to encode response: %v", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Debug("Added player to game", "playerId", playerId, "tableId", tableId)
}

func (a *RestApi) TogglePlayerReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	tableId := r.PathValue("tableId")
	playerId := r.PathValue("playerId")

	game, ok := a.Games[tableId]
	if !ok {
		slog.Error("Game not found", "tableId", tableId)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	playerIndex := -1
	for i, p := range game.Players {
		if p.Id == playerId {
			playerIndex = i
			break
		}
	}
	if playerIndex == -1 {
		slog.Error("Player not found", "playerId", playerId)
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	if game.State != blackjack.WaitingForPlayers {
		slog.Error("Game is not waiting for readiness", "tableId", tableId)
		http.Error(w, "Game is not waiting for readiness", http.StatusForbidden)
		return
	}

	game.TogglePlayerReady(playerId)
	player := game.Players[playerIndex]
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(player); err != nil {
		slog.Error(fmt.Sprintf("Failed to encode response: %v", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Debug("Toggled readiness for player", "playerId", playerId)
}

// nolint: cyclop
func (a *RestApi) PlayerAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	tableId := r.PathValue("tableId")
	playerId := r.PathValue("playerId")

	game, ok := a.Games[tableId]
	if !ok {
		slog.Error("Game not found", "tableId", tableId)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	playerIndex := -1
	for i, p := range game.Players {
		if p.Id == playerId {
			playerIndex = i
			break
		}
	}
	if playerIndex == -1 {
		slog.Error("Player not found", "playerId", playerId)
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	action := r.URL.Query().Get("action")

	switch action {
	case "hit":
		if !game.PlayerAction(playerIndex, blackjack.Hit) {
			http.Error(w, "Invalid action", http.StatusBadRequest)
			return
		}
		slog.Debug("Player hit", "playerId", playerId)
	case "stand":
		if !game.PlayerAction(playerIndex, blackjack.Stand) {
			http.Error(w, "Invalid action", http.StatusBadRequest)
			return
		}
		slog.Debug("Player stood", "playerId", playerId)
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}
}

func getRandomId() string {
	id, err := rand.Int(rand.Reader, big.NewInt(constant.MaxId))
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(id.Int64()))
}
