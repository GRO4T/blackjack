// TODO: Improve the file structure (types and functions are mixed up)
// TODO: Add API tests

package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/GRO4T/blackjack/deck"
)

const maxTableId = 1000
const maxPlayers = 1

// TODO: Think about whether chips and bets should be a part of Blackjack
type Game struct {
	Table   Blackjack
	TableId string
	Players []string
	Chips   []int
	Bets    []int
}

type GameDto struct {
	TableId string        `json:"tableId"`
	Players []string      `json:"players"`
	Hands   [][]deck.Card `json:"hands"`
	Chips   []int         `json:"chips"`
	Bets    []int         `json:"bets"`
}

// TODO: Search if there is more idiomatic way to do this
func buildDto(game Game) GameDto {
	return GameDto{
		TableId: game.TableId,
		Players: game.Players,
		Hands:   [][]deck.Card{game.Table.PlayerHand, game.Table.DealerHand},
		Chips:   game.Chips,
		Bets:    game.Bets,
	}
}

type Api struct {
	Games  map[string]Game
	Random *rand.Rand
}

func NewApi() Api {
	return Api{
		Games:  map[string]Game{},
		Random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (a *Api) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	log.Printf("POST /tables\n") // TODO: Can we simply make this automatic?

	tableId := strconv.Itoa(a.Random.Intn(maxTableId))
	newGame := Game{
		Table:   NewBlackjack(),
		TableId: tableId,
		Players: []string{},
		Chips:   []int{},
		Bets:    []int{},
	}
	newGame.Table.deal() // TODO: Maybe this should be in the initialization?
	a.Games[tableId] = newGame

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildDto(newGame))
	log.Printf("Created a new game %s\n", tableId)
}

func (a *Api) GetGameState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	log.Printf("GET /tables/%s\n", r.PathValue("tableId"))

	tableId := r.PathValue("tableId")
	game, ok := a.Games[tableId]
	if !ok {
		log.Printf("Game not found: %s\n", tableId)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildDto(game))
	log.Printf("Retrieved game state for %s\n", tableId)
}

func (a *Api) AddPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	log.Printf("POST /tables/players/%s\n", r.PathValue("tableId"))

	tableId := r.PathValue("tableId")
	game, ok := a.Games[tableId]
	if !ok {
		log.Printf("Game not found: %s\n", tableId)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	if len(game.Players) >= maxPlayers {
		log.Printf("Game is full: %s\n", tableId)
		http.Error(w, "Game is full", http.StatusForbidden)
		return
	}

	playerId := strconv.Itoa(a.Random.Intn(maxTableId))
	game.Players = append(game.Players, playerId)
	game.Chips = append(game.Chips, 100) // TODO: Remove magic number
	game.Bets = append(game.Bets, 0)
	a.Games[tableId] = game

	var resp struct {
		PlayerId string `json:"playerId"`
	}
	resp.PlayerId = playerId
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
	log.Printf("Added player %s to game %s\n", playerId, tableId)
}

func (a *Api) PlayerAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	log.Printf("POST /tables/%s/%s?action=%s\n", r.PathValue("tableId"), r.PathValue("playerId"), r.URL.Query().Get("action"))

	tableId := r.PathValue("tableId")
	game, ok := a.Games[tableId]
	if !ok {
		log.Printf("Game not found: %s\n", tableId)
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	playerId := r.PathValue("playerId")
	playerIndex := -1
	for i, p := range game.Players {
		if p == playerId {
			playerIndex = i
			break
		}
	}
	if playerIndex == -1 {
		log.Printf("Player not found: %s\n", playerId)
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	action := r.URL.Query().Get("action")

	// TODO: Add Hit and Stand functions to Blackjack
	if action == "hit" {
		game.Table.PlayerHand = append(game.Table.PlayerHand, game.Table.Deck[0])
		game.Table.Deck = game.Table.Deck[1:]
		game.Chips[playerIndex] -= game.Bets[playerIndex]
		game.Bets[playerIndex] += 10 // TODO: Allow player to control the bet amount
		log.Printf("Player %s hit\n", playerId)
	} else if action == "stand" {
		log.Printf("Player %s stood\n", playerId)
	} else {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildDto(game))
}
