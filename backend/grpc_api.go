package main

import (
	"context"

	pb "github.com/GRO4T/blackjack/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// nolint: gosec
func (s *BlackjackGrpcServer) GetGameState(c context.Context, r *pb.GetGameStateRequest) (*pb.GetGameStateResponse, error) {
	game, ok := s.Games[r.TableId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Game not found")
	}

	pbHands := []*pb.Hand{}
	for _, cards := range game.Hands {
		pbCards := []*pb.Card{}
		for _, card := range cards {
			pbCards = append(pbCards, &pb.Card{
				Rank: int32(card.Rank), // nolint: gosec
				Suit: int32(card.Suit), // nolint: gosec
			})
		}
		pbHands = append(pbHands, &pb.Hand{Cards: pbCards})
	}

	pbPlayers := []*pb.Player{}
	for _, player := range game.Players {
		pbPlayers = append(pbPlayers, &pb.Player{
			Name:    player.Name,
			IsReady: player.IsReady,
			Chips:   int32(player.Chips),
			Bet:     int32(player.Bet), //nolint: gosec
			Outcome: pb.Outcome(player.Outcome),
		})
	}

	return &pb.GetGameStateResponse{
		Hands:         pbHands,
		Players:       pbPlayers,
		State:         pb.State(game.State),
		CurrentPlayer: int32(game.CurrentPlayer),
	}, nil
}

func (s *BlackjackGrpcServer) AddPlayer(c context.Context, r *pb.AddPlayerRequest) (*pb.AddPlayerResponse, error) {
	game, ok := s.Games[r.TableId]

	if !ok {
		return nil, status.Errorf(codes.NotFound, "Game not found")
	}
	if len(game.Players) >= maxPlayers {
		return nil, status.Errorf(codes.FailedPrecondition, "Game is full")
	}
	if game.State != WaitingForPlayers {
		return nil, status.Errorf(codes.FailedPrecondition, "Game has already started")
	}

	playerId := getRandomId()
	game.AddPlayer(playerId, "Bob")
	return &pb.AddPlayerResponse{PlayerId: playerId}, nil
}

// nolint: gosec
func (s *BlackjackGrpcServer) TogglePlayerReady(c context.Context, r *pb.TogglePlayerReadyRequest) (*pb.Player, error) {
	game, ok := s.Games[r.TableId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Game not found")
	}

	playerIndex := -1
	for i, p := range game.Players { // TODO(refactor): Change players to map
		if p.Id == r.PlayerId {
			playerIndex = i
			break
		}
	}
	if playerIndex == -1 {
		return nil, status.Errorf(codes.NotFound, "Player not found")
	}

	if game.State != WaitingForPlayers {
		return nil, status.Errorf(codes.FailedPrecondition, "Game is not waiting for readiness")
	}

	game.TogglePlayerReady(r.PlayerId)
	player := game.Players[playerIndex]
	return &pb.Player{
		Name:    player.Name,
		IsReady: player.IsReady,
		Chips:   int32(player.Chips),
		Bet:     int32(player.Bet),
		Outcome: pb.Outcome(player.Outcome),
	}, nil
}

func (s *BlackjackGrpcServer) PlayerAction(c context.Context, r *pb.PlayerActionRequest) (*emptypb.Empty, error) {
	game, ok := s.Games[r.TableId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Game not found")
	}

	playerIndex := -1
	for i, p := range game.Players {
		if p.Id == r.PlayerId {
			playerIndex = i
			break
		}
	}
	if playerIndex == -1 {
		return nil, status.Errorf(codes.NotFound, "Player not found")
	}

	switch r.Action {
	case pb.Action_HIT:
		if !game.PlayerAction(playerIndex, Hit) {
			return nil, status.Errorf(codes.FailedPrecondition, "Invalid action")
		}
	case pb.Action_STAND:
		if !game.PlayerAction(playerIndex, Stand) {
			return nil, status.Errorf(codes.FailedPrecondition, "Invalid action")
		}
	}

	return &emptypb.Empty{}, nil
}
