// TODO: Document the API using Swagger
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	pb "github.com/GRO4T/blackjack/grpc"
	"google.golang.org/grpc"
)

const (
	ServerAddr = "localhost:8080"
)

func grpcServer() {
	listener, err := net.Listen("tcp", ServerAddr)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to listen: %v", err))
	}
	s := grpc.NewServer()
	pb.RegisterBlackjackServer(s, NewGrpcServer())
	slog.Info(fmt.Sprintf("Starting gRPC server on %s", ServerAddr))
	if err := s.Serve(listener); err != nil {
		slog.Error(fmt.Sprintf("Failed to serve: %v", err))
	}
}

// nolint: mnd
func restApiServer() {
	api := NewRestApi()
	mux := http.NewServeMux()
	mux.HandleFunc("/tables", api.CreateGame)
	mux.HandleFunc("/tables/{tableId}", api.GetGameState)
	mux.HandleFunc("/tables/ready/{tableId}/{playerId}", api.TogglePlayerReady)
	mux.HandleFunc("/tables/players/{tableId}", api.AddPlayer)
	mux.HandleFunc("/tables/{tableId}/{playerId}", api.PlayerAction)
	s := &http.Server{
		Addr:           ServerAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	slog.Info(fmt.Sprintf("Starting REST server on %s", ServerAddr))
	slog.Error(s.ListenAndServe().Error())
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	isGrpc := flag.Bool("grpc", false, "Start gRPC server instead of REST")
	flag.Parse()
	if *isGrpc {
		grpcServer()
	} else {
		restApiServer()
	}
}
