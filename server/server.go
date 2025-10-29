package main

import (
	pb "ChitChat/proto"
	"sync"
)

func main() {
	// Empty func is empty...
}

// Server constructor
type chitChatServer struct {
	pb.UnimplementedChitChatServer

	mu      sync.Mutex  // Protects the server's latest message from being overridden twice at once
	lastMsg *pb.Message // The server's latest message
}
