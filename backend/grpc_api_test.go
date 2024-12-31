package main

import (
	"context"
	"log"
	"net"
	"testing"

	pb "github.com/GRO4T/blackjack/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// nolint: ireturn
func Setup_TestGrpcServer(t *testing.T) (*BlackjackGrpcServer, pb.BlackjackClient) {
	t.Helper()
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	serviceRegistrar := grpc.NewServer()
	t.Cleanup(func() {
		serviceRegistrar.Stop()
	})
	server := NewGrpcServer()
	pb.RegisterBlackjackServer(serviceRegistrar, server)

	go func() {
		if err := serviceRegistrar.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	t.Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		t.Fatal(err)
	}

	client := pb.NewBlackjackClient(conn)
	return server, client
}

func TestGrpcServer_CreateGame(t *testing.T) {
	// Arrange
	_, client := Setup_TestGrpcServer(t)

	// Act
	res, err := client.CreateGame(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if res.TableId == "" {
		t.Fatal("TableId is empty")
	}
}

func TestGrpcServer_GetGameState(t *testing.T) {
	// Arrange
	server, client := Setup_TestGrpcServer(t)
	game := NewBlackjack()
	game.AddPlayer("1", "Player 1")
	server.Games["1"] = &game

	// Act
	res, err := client.GetGameState(context.Background(), &pb.GetGameStateRequest{TableId: "1"})
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if len(res.Players) != 1 {
		t.Fatalf("Expected 1 player; got %v", len(res.Players))
	}
}

func TestGrpcServer_AddPlayer(t *testing.T) {
	// Arrange
	server, client := Setup_TestGrpcServer(t)
	game := NewBlackjack()
	server.Games["1"] = &game

	// Act
	_, err := client.AddPlayer(context.Background(), &pb.AddPlayerRequest{TableId: "1"})
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if len(server.Games["1"].Players) != 1 {
		t.Errorf("Expected 1 player; got %v", len(server.Games["1"].Players))
	}
}

func TestGrpcServer_TogglePlayerReadyWhenPlayerNotReady(t *testing.T) {
	// Arrange
	server, client := Setup_TestGrpcServer(t)
	game := NewBlackjack()
	game.AddPlayer("1", "Player 1")
	server.Games["1"] = &game

	// Act
	_, err := client.TogglePlayerReady(
		context.Background(),
		&pb.TogglePlayerReadyRequest{TableId: "1", PlayerId: "1"},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if !server.Games["1"].Players[0].IsReady {
		t.Error("Player is not ready")
	}
	if server.Games["1"].State != CardsDealt {
		t.Error("Game is not in CardsDealt state")
	}
}

func TestGrpcServer_TogglePlayerReadyWhenPlayerReady(t *testing.T) {
	// Arrange
	server, client := Setup_TestGrpcServer(t)
	game := NewBlackjack()
	game.AddPlayer("1", "Player 1")
	server.Games["1"] = &game
	server.Games["1"].Players[0].IsReady = true

	// Act
	_, err := client.TogglePlayerReady(
		context.Background(),
		&pb.TogglePlayerReadyRequest{TableId: "1", PlayerId: "1"},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if server.Games["1"].Players[0].IsReady {
		t.Error("Player is ready")
	}
}

func TestGrpcServer_PlayerAction(t *testing.T) {
	// Arrange
	server, client := Setup_TestGrpcServer(t)
	game := NewBlackjack()
	game.AddPlayer("1", "Player 1")
	game.Deal()
	game.State = CardsDealt
	server.Games["1"] = &game

	// Act
	_, err := client.PlayerAction(
		context.Background(),
		&pb.PlayerActionRequest{TableId: "1", PlayerId: "1", Action: pb.Action_HIT},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if len(server.Games["1"].GetPlayerHand(0)) != 3 {
		t.Errorf("Expected 3 cards; got %v", len(server.Games["1"].GetPlayerHand(0)))
	}
}

func TestGrpcServer_SimpleGame(t *testing.T) {
	_, client := Setup_TestGrpcServer(t)
	ctx := context.Background()

	// Create game
	createGameResp, err := client.CreateGame(ctx, &emptypb.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	tableId := createGameResp.TableId

	// Add player
	addPlayerResp, err := client.AddPlayer(ctx, &pb.AddPlayerRequest{TableId: tableId})
	if err != nil {
		t.Fatal(err)
	}
	playerId := addPlayerResp.PlayerId

	// Toggle player ready
	_, err = client.TogglePlayerReady(ctx, &pb.TogglePlayerReadyRequest{TableId: tableId, PlayerId: playerId})
	if err != nil {
		t.Fatal(err)
	}

	// Player hit
	_, err = client.PlayerAction(
		ctx,
		&pb.PlayerActionRequest{TableId: tableId, PlayerId: playerId, Action: pb.Action_HIT},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Check game outcome
	gameState, err := client.GetGameState(ctx, &pb.GetGameStateRequest{TableId: tableId})
	if err != nil {
		t.Fatal(err)
	}
	if gameState.Players[0].Outcome == pb.Outcome_UNDECIDED {
		t.Error("Expected player outcome to be decided")
	}
}
