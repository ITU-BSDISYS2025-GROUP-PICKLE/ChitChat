package main

import (
	pb "ChitChat/proto"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
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

	// Set a client name (username)
	println()
	print("Please enter your name (public to other users): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	clientName := scanner.Text()
	println()

	// Read the latest message from the server (to get timestamp)
	in, err := stream.Recv()
	if err != nil {
		log.Fatalf("client.RouteChat failed: %v", err)
	}

	// Create Lamport timestamp
	var lamportTime int64 = in.Time + 1

	// Send arrival message to server
	arrival := fmt.Sprintf("Participant %s joined the ChitChat at logical time %d", clientName, lamportTime)
	err = stream.Send(&pb.Message{Message: arrival, ClientName: clientName, Time: lamportTime})
	if err != nil {
		log.Fatalf("client.RouteChat: stream.Send(%v) failed: %v", err, err)
	}

	// Continually read latest message from the server
	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				log.Fatalf("client.RouteChat failed: %v", err)
			}

			// This logic ensures only new messages received from the server are printed
			if in.Time >= lamportTime {
				fmt.Printf("T%d | %s > %s\n", in.Time, in.ClientName, in.Message)
				lamportTime = in.Time + 1
			}
		}
	}()

	// Continually send input messages to the server
	for {
		scanner.Scan()
		input := scanner.Text()

		// Ignore empty input
		if strings.Trim(input, " ") == "" {
			continue
		}

		if len(input) > 127 {
			fmt.Println("Error: message too long")
		}

		// Disconnet if input is '/exit'
		if input == "/exit" {
			break
		}

		if err := stream.Send(&pb.Message{Message: input, ClientName: clientName, Time: lamportTime}); err != nil {
			log.Fatalf("client.RouteChat: stream.Send(%v) failed: %v", err, err)
		}
	}

	// Send departure message to server
	departure := fmt.Sprintf("Participant %s left the ChitChat at logical time %d", clientName, lamportTime)
	err = stream.Send(&pb.Message{Message: departure, ClientName: clientName, Time: lamportTime})
	if err != nil {
		log.Fatalf("client.RouteChat: stream.Send(%v) failed: %v", err, err)
	}

	// Receive message from server so this client sees the departure message go through
	in, err = stream.Recv()
	if err != nil {
		log.Fatalf("client.RouteChat failed: %v", err)
	}
	if in.Time >= lamportTime {
		fmt.Printf("T%d | %s > %s\n", in.Time, in.ClientName, in.Message)
		lamportTime = in.Time + 1
	}

	// Close the receiving stream, let program terminate
	stream.CloseSend()
}
