package rest

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/GRO4T/bjack-api/blackjack"
	"github.com/GRO4T/bjack-api/constant"
	"github.com/gorilla/websocket"
)

type RestApi struct {
	Games      map[string]*blackjack.Blackjack
	Websockets map[string][]*websocket.Conn // TODO: Test if the websockets will close automatically when the server is killed.
}

type CreateGameRequest struct {
	PlayerName string `json:"playerName"`
}

type CreateGameResponse struct {
	TableId string `json:"tableId"`
}

type AddPlayerRequest struct {
	PlayerName string `json:"playerName"`
}

type AddPlayerResponse struct {
	PlayerId string `json:"playerId"`
}

func NewApi() RestApi {
	return RestApi{
		Games:      map[string]*blackjack.Blackjack{},
		Websockets: map[string][]*websocket.Conn{},
	}
}

func (a *RestApi) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var reqData CreateGameRequest
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to decode request: %v", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if reqData.PlayerName == "" {
		http.Error(w, "Player name cannot be empty", http.StatusBadRequest)
		return
	}

	tableId := getRandomId()
	newGame := blackjack.New(func() {
		for _, ws := range a.Websockets[tableId] {
			if err := ws.WriteMessage(websocket.TextMessage, []byte("NewState")); err != nil {
				slog.Error(fmt.Sprintf("Failed to send data via websocket: %v", err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
	})
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
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	var reqData AddPlayerRequest
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to decode request: %v", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if reqData.PlayerName == "" {
		http.Error(w, "Player name cannot be empty", http.StatusBadRequest)
		return
	}

	newPlayer, err := game.AddPlayer(reqData.PlayerName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp AddPlayerResponse
	resp.PlayerId = newPlayer.Id
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error(fmt.Sprintf("Failed to encode response: %v", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Debug("Added player to game", "playerId", newPlayer.Id, "tableId", tableId)
}

func (a *RestApi) RemovePlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	tableId := r.PathValue("tableId")
	playerId := r.PathValue("playerId")

	game, ok := a.Games[tableId]
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	if err := game.RemovePlayer(playerId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Debug("Removed player from game", "playerId", playerId, "tableId", tableId)
}

func (a *RestApi) TogglePlayerReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	tableId := r.PathValue("tableId")
	playerId := r.PathValue("playerId")

	game, ok := a.Games[tableId]
	if !ok {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	player, err := game.TogglePlayerReady(playerId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to toggle readiness: %v", err), http.StatusBadRequest)
		return
	}

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
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	action := r.URL.Query().Get("action")

	switch action {
	case "hit":
		if err := game.PlayerAction(playerId, blackjack.Hit); err != nil {
			http.Error(w, fmt.Sprintf("Invalid action: %v", err), http.StatusInternalServerError)
			return
		}
		slog.Debug("Player hit", "playerId", playerId)
	case "stand":
		if err := game.PlayerAction(playerId, blackjack.Stand); err != nil {
			http.Error(w, fmt.Sprintf("Invalid action: %v", err), http.StatusInternalServerError)
			return
		}
		slog.Debug("Player stood", "playerId", playerId)
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Accepting all requests
	},
}

func (a *RestApi) AddStateObserver(w http.ResponseWriter, r *http.Request) {
	tableId := r.PathValue("tableId")

	ws, _ := upgrader.Upgrade(w, r, nil)
	gameSubscribers, ok := a.Websockets[tableId]
	if ok {
		a.Websockets[tableId] = append(gameSubscribers, ws)
	} else {
		a.Websockets[tableId] = []*websocket.Conn{ws}
	}
}

func getRandomId() string {
	id, err := rand.Int(rand.Reader, big.NewInt(constant.MaxId))
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(id.Int64()))
}
