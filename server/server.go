package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	pb "auction/auction"

	"google.golang.org/grpc"
)

type AuctionState struct {
	HighestBid    int32
	HighestBidder string
}

type AuctionServer struct {
	pb.UnimplementedAuctionServer
	serverNodes   map[int32]pb.AuctionClient
	mu            sync.Mutex
	isLeader      bool
	nodeID        int32
	state         AuctionState
	auctionActive bool
	bidTimeout    time.Duration
	bidEndTime    time.Time
}

var ports = []string{":5050", ":5051", ":5052", ":5053"}

func NewAuctionServer(nodeID int32, isLeader bool, bidTimeout time.Duration) *AuctionServer {
	return &AuctionServer{
		serverNodes:   make(map[int32]pb.AuctionClient),
		isLeader:      isLeader,
		nodeID:        nodeID,
		state:         AuctionState{HighestBid: 0, HighestBidder: ""},
		auctionActive: true,
		bidTimeout:    bidTimeout,
		bidEndTime:    time.Now().Add(bidTimeout),
	}
}

func startServer(port string, nodeID int32, isLeader bool, bidTimeout time.Duration) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create and register the AuctionServer
	server := NewAuctionServer(nodeID, isLeader, bidTimeout)
	grpcServer := grpc.NewServer()
	pb.RegisterAuctionServer(grpcServer, server)

	// Start the auction timer
	go server.runAuction()

	log.Printf("Server %d listening on %v", nodeID, lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *AuctionServer) runAuction() {
	for {
		time.Sleep(time.Second)

		// Check if auction has ended
		s.mu.Lock()
		if time.Now().After(s.bidEndTime) && s.auctionActive {
			s.auctionActive = false
			log.Printf("Auction ended. Winner: %s with bid %d", s.state.HighestBidder, s.state.HighestBid)
		}
		s.mu.Unlock()

		// If the auction is over, exit the loop
		if !s.auctionActive {
			break
		}
	}
}

// Bid handling
func (s *AuctionServer) Bid(ctx context.Context, amount *pb.Amount) (*pb.Ack, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.auctionActive || time.Now().After(s.bidEndTime) {
		return &pb.Ack{Message: "auction over"}, nil
	}

	if amount.Amount <= s.state.HighestBid {
		return &pb.Ack{Message: "fail"}, nil
	}

	// If the server is the leader, propagate the bid to other nodes
	if s.isLeader {
		success := s.propagateBid(amount.Amount)
		if !success {
			return &pb.Ack{Message: "exception"}, nil
		}
	}

	// Update auction state
	s.state.HighestBid = amount.Amount
	s.state.HighestBidder = "Bidder_ID"
	return &pb.Ack{Message: "success"}, nil
}

// Result handling
func (s *AuctionServer) Result(ctx context.Context, void *pb.Void) (*pb.Outcome, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.auctionActive && time.Now().Before(s.bidEndTime) {
		return &pb.Outcome{
			HighestBid:    s.state.HighestBid,
			HighestBidder: "in progress",
		}, nil
	}

	// Auction is over, return final outcome
	return &pb.Outcome{
		HighestBid:    s.state.HighestBid,
		HighestBidder: s.state.HighestBidder,
	}, nil
}

// Propagate bid to other server nodes
func (s *AuctionServer) propagateBid(amount int32) bool {
	ackCount := 0
	for _, client := range s.serverNodes {
		ack, err := client.Bid(context.Background(), &pb.Amount{Amount: amount})
		if err == nil && ack.Message == "success" {
			ackCount++
		}
	}

	// Majority of nodes must acknowledge the bid
	return ackCount >= len(s.serverNodes)/2
}

// Leader election mechanism (simplified)
func (s *AuctionServer) electNewLeader() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If the server is the current leader, then it should step down
	if s.isLeader {
		s.isLeader = false
		log.Printf("Server %d stepped down as leader", s.nodeID)
	}

	// Election logic - make the node with the lowest ID the new leader
	s.isLeader = true
	log.Printf("Server %d is now the leader", s.nodeID)
}

func main() {
	// Start servers on predefined ports
	for i, port := range ports {
		isLeader := (i == 0) // First server is the initial leader
		go startServer(port, int32(i), isLeader, 100*time.Second)
	}

	// Block forever (let servers run)
	select {}
}
