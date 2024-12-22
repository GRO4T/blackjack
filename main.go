package main

import (
	"log"
	"net/http"
)

func main() {
	api := NewApi()
	http.HandleFunc("/tables", api.CreateGame)
	http.HandleFunc("/tables/{tableId}", api.GetGameState)
	http.HandleFunc("/tables/ready/{tableId}/{playerId}", api.TogglePlayerReady)
	http.HandleFunc("/tables/players/{tableId}", api.AddPlayer)
	http.HandleFunc("/tables/{tableId}/{playerId}", api.PlayerAction)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
