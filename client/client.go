package main

import (
	pb "ChitChat/proto"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// Create connection
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close() // Defer closing the connection until main func is done running

	// Create client
	client := pb.NewChitChatClient(conn)

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// Get stream via RPC
	stream, err := client.ChatRoom(ctx)
	if err != nil {
		log.Fatalf("client.RouteChat failed: %v", err)
	}
}
