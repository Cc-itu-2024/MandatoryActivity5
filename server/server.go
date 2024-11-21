package main

import (
	"context"
	"log"
	"net"
	"os"
	"sync"
	"time"

	pb "auction/auction"

	"google.golang.org/grpc"
)

var ports = []string{":5050", ":5051", ":5052",  ":5053"}
var Auctionactive bool

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

type AuctionState struct {
	HighestBid    int32  
	HighestBidder string 
}


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

func main() {
	logFile, err := os.OpenFile("auction.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	// Start servers on predefined ports
	for i, port := range ports {
		isLeader := (i == 0) // First server is the initial leader
		go startServer(port, int32(i), isLeader, 100*time.Second)
	}

	// Block forever
	select {}
}



func (s *AuctionServer) runAuction() {
	for {
		time.Sleep(time.Second) // Check periodically

		s.mu.Lock()
		if time.Now().After(s.bidEndTime) && s.auctionActive {
			s.auctionActive = false
			log.Printf("Auction ended. Winner: %s with bid %d", s.state.HighestBidder, s.state.HighestBid)
		}
		s.mu.Unlock()

		// Exit the loop if the auction is over
		if !s.auctionActive {
			break
		}
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

	go server.runAuction() // Start the auction timer

	log.Printf("Server %d listening on %v", nodeID, lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}



func (s *AuctionServer) Bid(ctx context.Context, amount *pb.Amount) (*pb.Ack, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.auctionActive || time.Now().After(s.bidEndTime) {
		return &pb.Ack{Message: "auction over"}, nil
	}

	if amount.Amount <= s.state.HighestBid {
		return &pb.Ack{Message: "fail"}, nil
	}

	if s.isLeader {
		success := s.propagateBid(amount.Amount)
		if !success {
			return &pb.Ack{Message: "exception"}, nil
		}
	}

	s.state.HighestBid = amount.Amount
	s.state.HighestBidder = "Bidder_ID" 
	return &pb.Ack{Message: "success"}, nil
}


func ack(){

}

func outcome(){

}

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

func (s *AuctionServer) propagateBid(amount int32) bool {
	ackCount := 0
	for _, client := range s.serverNodes {
		ack, err := client.Bid(context.Background(), &pb.Amount{Amount: amount})
		if err == nil && ack.Message == "success" {
			ackCount++
		}
	}

	
	return ackCount >= len(s.serverNodes)/2
}

func (s *AuctionServer) electNewLeader(){


}


