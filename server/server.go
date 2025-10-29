package main

import (
	pb "ChitChat/proto"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"google.golang.org/grpc"
)

// Server constructor
type chitChatServer struct {
	pb.UnimplementedChitChatServer

	mu      sync.Mutex  // Protects the server's latest message from being overridden twice at once
	lastMsg *pb.Message // The server's latest message
}

func main() {
	flag.Parse()

	// Create listener
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create server
	grpcServer := grpc.NewServer()
	chitChatServer := &chitChatServer{lastMsg: &pb.Message{ClientName: "", Message: "", Time: 0}}
	pb.RegisterChitChatServer(grpcServer, chitChatServer)

	// Create log file
	year := fmt.Sprint(time.Now().Year())
	month := fmt.Sprint(int(time.Now().Month()))
	day := fmt.Sprint(time.Now().Day())
	hour := fmt.Sprint(time.Now().Hour())
	minute := fmt.Sprint(time.Now().Minute())
	second := fmt.Sprint(time.Now().Second())

	fileName := "server/logs/log-" + year + month + day + "-" + hour + minute + second + ".txt"

	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	// Start server
	fmt.Println("Server listening on: localhost:50051")
	log.Println("Server started. Listening on: localhost:50051")
	defer log.Println("Server shutting down.")
	grpcServer.Serve(lis)
}
