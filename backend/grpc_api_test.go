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

func TestGrpcServer_CreateGame(t *testing.T) {
	// Arrange
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	s := grpc.NewServer()
	t.Cleanup(func() {
		s.Stop()
	})
	pb.RegisterBlackjackServer(s, NewGrpcServer())

	go func() {
		if err := s.Serve(lis); err != nil {
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

	// Act
	client := pb.NewBlackjackClient(conn)
	res, err := client.CreateGame(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	if res.TableId == "" {
		t.Fatal("TableId is empty")
	}
}
