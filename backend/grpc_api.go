package main

import (
	"context"

	pb "github.com/GRO4T/blackjack/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type BlackjackGrpcServer struct {
	pb.UnimplementedBlackjackServer
	Games map[string]*Blackjack
}

func NewGrpcServer() *BlackjackGrpcServer {
	return &BlackjackGrpcServer{
		Games: map[string]*Blackjack{},
	}
}

func (s *BlackjackGrpcServer) CreateGame(context.Context, *emptypb.Empty) (*pb.CreateGameResponse, error) {
	tableId := getRandomId()
	newGame := NewBlackjack()
	s.Games[tableId] = &newGame
	return &pb.CreateGameResponse{TableId: tableId}, nil
}
