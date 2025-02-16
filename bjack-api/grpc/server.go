package grpc

import (
	"context"
	"crypto/rand"
	"math/big"
	"strconv"

	"github.com/GRO4T/bjack-api/blackjack"
	"github.com/GRO4T/bjack-api/constant"
	pb "github.com/GRO4T/bjack-api/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type BlackjackServer struct {
	pb.UnimplementedBlackjackServer
	Games map[string]*blackjack.Blackjack
}

func NewServer() *BlackjackServer {
	return &BlackjackServer{
		Games: map[string]*blackjack.Blackjack{},
	}
}

func (s *BlackjackServer) CreateGame(context.Context, *emptypb.Empty) (*pb.CreateGameResponse, error) {
	tableId := getRandomId()
	newGame := blackjack.New(nil)
	s.Games[tableId] = &newGame
	return &pb.CreateGameResponse{TableId: tableId}, nil
}

// nolint: gosec
func (s *BlackjackServer) GetGameState(c context.Context, r *pb.GetGameStateRequest) (*pb.GetGameStateResponse, error) {
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

func (s *BlackjackServer) AddPlayer(c context.Context, r *pb.AddPlayerRequest) (*pb.AddPlayerResponse, error) {
	game, ok := s.Games[r.TableId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Game not found")
	}
	newPlayer, err := game.AddPlayer("Bob")
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "Failed to add player: %v", err)
	}
	return &pb.AddPlayerResponse{PlayerId: newPlayer.Id}, nil
}

// nolint: gosec
func (s *BlackjackServer) TogglePlayerReady(c context.Context, r *pb.TogglePlayerReadyRequest) (*pb.Player, error) {
	game, ok := s.Games[r.TableId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Game not found")
	}

	player, err := game.TogglePlayerReady(r.PlayerId)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "Failed to toggle readiness: %v", err)
	}

	return &pb.Player{
		Name:    player.Name,
		IsReady: player.IsReady,
		Chips:   int32(player.Chips),
		Bet:     int32(player.Bet),
		Outcome: pb.Outcome(player.Outcome),
	}, nil
}

func (s *BlackjackServer) PlayerAction(c context.Context, r *pb.PlayerActionRequest) (*emptypb.Empty, error) {
	game, ok := s.Games[r.TableId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Game not found")
	}

	switch r.Action {
	case pb.Action_HIT:
		if err := game.PlayerAction(r.PlayerId, blackjack.Hit); err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "Invalid action")
		}
	case pb.Action_STAND:
		if err := game.PlayerAction(r.PlayerId, blackjack.Stand); err != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "Invalid action")
		}
	}

	return &emptypb.Empty{}, nil
}

func getRandomId() string {
	id, err := rand.Int(rand.Reader, big.NewInt(constant.MaxId))
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(id.Int64()))
}
