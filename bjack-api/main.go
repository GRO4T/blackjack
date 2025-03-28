// TODO: Document the API using Swagger
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	bgrpc "github.com/GRO4T/bjack-api/grpc"
	pb "github.com/GRO4T/bjack-api/proto"
	"github.com/GRO4T/bjack-api/rest"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

const (
	ServerAddr = "0.0.0.0:8000"
)

func grpcServer() {
	listener, err := net.Listen("tcp", ServerAddr) //nolint:gosec
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to listen: %v", err))
	}
	s := grpc.NewServer()
	pb.RegisterBlackjackServer(s, bgrpc.NewServer())
	slog.Info(fmt.Sprintf("Starting gRPC server on %s", ServerAddr))
	if err := s.Serve(listener); err != nil {
		slog.Error(fmt.Sprintf("Failed to serve: %v", err))
	}
}

// nolint: mnd
func restApiServer() {
	api := rest.NewApi()

	mux := http.NewServeMux()
	mux.HandleFunc("/tables", api.CreateGame)
	mux.HandleFunc("/tables/{tableId}", api.GetGameState)
	mux.HandleFunc("/tables/ready/{tableId}/{playerId}", api.TogglePlayerReady)
	mux.HandleFunc("/tables/players/{tableId}", api.AddPlayer)
	mux.HandleFunc("/tables/players/{tableId}/{playerId}", api.RemovePlayer)
	mux.HandleFunc("/tables/{tableId}/{playerId}", api.PlayerAction)
	mux.HandleFunc("/state-updates/{tableId}", api.AddStateObserver)

	uiUrl, ok := os.LookupEnv("UI_URL")
	if !ok {
		slog.Error("UI_URL not provided")
		os.Exit(1)
	}
	corsMux := cors.New(cors.Options{
		AllowedOrigins: []string{uiUrl},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler(mux)

	s := &http.Server{
		Addr:           ServerAddr,
		Handler:        corsMux,
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
