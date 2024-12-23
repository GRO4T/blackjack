package main

import (
	"log/slog"
	"net/http"
	"time"
)

// nolint: mnd
func main() {
	api := NewApi()
	mux := http.NewServeMux()
	mux.HandleFunc("/tables", api.CreateGame)
	mux.HandleFunc("/tables/{tableId}", api.GetGameState)
	mux.HandleFunc("/tables/ready/{tableId}/{playerId}", api.TogglePlayerReady)
	mux.HandleFunc("/tables/players/{tableId}", api.AddPlayer)
	mux.HandleFunc("/tables/{tableId}/{playerId}", api.PlayerAction)
	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	slog.Info("Starting server on :8080\n")
	slog.Error(s.ListenAndServe().Error())
}
