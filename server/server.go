package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "auction/auction"

	"google.golang.org/grpc"
)

var (
	nodeId    = flag.Int("nodeId", -1, "The nodeId")
	maxNode = flag.Int("maxNode", -1, "The maximum number of nodes")
	basePort  = flag.Int("basePort", 50000, "The base port number. The server will listen on this port + nodeId")
	min       = flag.Int("min", 1, "The durition of the auction in minutes")
)

var AuctionOngoing = false

type server struct {
	pb.UnimplementedAuctionServer

	// Auction state
	auctionState    pb.AUCTION_STATUS
	highestBid      int32
	highestBidderId int32

	// Node state
	nodeId   int
	leaderId int

	// Mutex to lock the server
	mu sync.Mutex
}

func (s *server) Bid(ctx context.Context, req *pb.BidRequest) (*pb.BidResponse, error) {
	if s.leaderId == s.nodeId {
		// I'm the leader. I can accept the bid
		s.mu.Lock()
		defer s.mu.Unlock()

		// TODO Check if the auction is ongoing
		// If finished, return NotAccepted/Finished
		if (s.auctionState == pb.AUCTION_STATUS_FINISHED) {
			return &pb.BidResponse{Status: pb.BID_STATUS_ISFINISHED, Bid: s.highestBid, BidderId: s.highestBidderId}, nil
		}

		if req.Bid > s.highestBid {
			s.highestBid = req.Bid
			s.highestBidderId = req.BidderId
			log.Printf("Node %d: New highest bid of %d by bidder %d", s.nodeId, s.highestBid, s.highestBidderId)

			// Now broadcast the state to all nodes with lower id
			s.broadcastState()

			return &pb.BidResponse{Status: pb.BID_STATUS_ACCEPTED, Bid: s.highestBid, BidderId: s.highestBidderId}, nil
		}

		return &pb.BidResponse{Status: pb.BID_STATUS_NOTACCEPTED, Bid: s.highestBid, BidderId: s.highestBidderId}, nil
	}

	client, conn, _ := s.getLeaderAuctionClient()

	if client == nil {
		return s.Bid(ctx, req) // I am the leader now. Retry the bid. Call myself
	}

	defer conn.Close()
	return client.Bid(ctx, req)
}

// Returns auction state.
// If the current node is the leader, it returns the state. Otherwise, it forwards the request to the leader
func (s *server) Result(ctx context.Context, req *pb.Empty) (*pb.ResultResponse, error) {
	if s.leaderId == *nodeId {
		// I'm the leader. I can return the result

		// Lock the server
		s.mu.Lock()
		defer s.mu.Unlock()

		return &pb.ResultResponse{Status: s.auctionState, HighestBid: s.highestBid, BidderId: s.highestBidderId}, nil
	}

	client, conn, _ := s.getLeaderAuctionClient()

	if client == nil {
		return s.Result(ctx, req) // I am the leader now. Retry the result
	}

	defer conn.Close()

	return client.Result(ctx, req)
}

// GetLeaderAuctionClient returns the client to the leader node or nil if the current node is the leader
// If the current node is not the leader or the existing leader doens't reply, it starts the election process
// If the current node becomes the leader after the election process, it will return nil and the caller should retry the operation
func (s *server) getLeaderAuctionClient() (pb.AuctionClient, *grpc.ClientConn, error) {
	if s.leaderId == -1 {
		s.startElection()
	}

	if s.leaderId == s.nodeId {
		return nil, nil, nil
	}

	node := fmt.Sprintf("localhost:%d", *basePort+s.leaderId)
	conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second*2))
	if err != nil {
		log.Printf("Node %d: Could not connect to node %s", s.nodeId, node)

		s.startElection()
		return s.getLeaderAuctionClient()
	}

	return pb.NewAuctionClient(conn), conn, nil
}

// StartElection starts the election process
// It sends an election request to all nodes with higher id
func (s *server) startElection() {
	log.Printf("Node %d: Starting election", s.nodeId)

	for n := s.nodeId + 1; n < *maxNode; n++ {
		node := fmt.Sprintf("localhost:%d", *basePort+n)
		log.Printf("Node %d: Try to contact NodeId: %d on address: %s", s.nodeId, n, node)

		conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
		if err != nil {
			log.Printf("Node %d: Could not connect to node %s", s.nodeId, node)
			continue
		}
		defer conn.Close()

		client := pb.NewAuctionClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		log.Printf("Node %d: SendElection request to node: %d", s.nodeId, n)
		r, err := client.SendElection(ctx, &pb.ElectionRequest{NodeId: int32(s.nodeId)})
		if err != nil {
			log.Printf("Node %d: No response from node %s", s.nodeId, err)
			continue
		}

		if !r.Ok {
			// If we accidentially sent an election request to a node with lower id, we should ignore it
			log.Printf("Node %d: Received not Ok from node %s. Lower id than this node.", s.nodeId, node)
			continue
		}

		// We received an OK from a node with higher id. We are not the leader. Exit the election process
		log.Printf("Node %d: Received OK from node %s", s.nodeId, node)
		return
	}

	// No response from any node. I am the leader
	// Lock the server and update the leader
	s.mu.Lock()
	defer s.mu.Unlock()

	s.leaderId = s.nodeId
	log.Printf("Node %d: I am the leader", s.nodeId)

	// Announce the new leader (this) to all nodes with lower id
	s.announceLeader()
}

// AnnounceLeader announces the leader to all nodes with lower id
// It sends an AnnounceLeader request to all nodes with lower id
func (s *server) announceLeader() {
	log.Printf("Node %d: Announcing leader to nodes with lower id.", s.nodeId)
	for n := 0; n < s.nodeId; n++ {
		node := fmt.Sprintf("localhost:%d", *basePort+n)

		conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
		if err != nil {
			log.Printf("Node %d: Could not connect to node %s", s.nodeId, node)
			continue
		}
		defer conn.Close()

		c := pb.NewAuctionClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err = c.AnnounceLeader(ctx, &pb.AnnounceLeaderRequest{LeaderId: int32(s.nodeId)})
		if err != nil {
			log.Printf("Node %d: Could not announce leader to node %s", s.nodeId, node)
			continue
		}

		log.Printf("Node %d: Announced leader to node %s", s.nodeId, node)
	}

	log.Printf("Node %d: Announced leader to all nodes", s.nodeId)
}

// BroadcastState broadcasts the state to all nodes with lower id
// It sends a SyncState request to all nodes with lower id to sync the state to them
// This is used to keep all nodes in sync with the leader in case of leader change
// The leader calls this to replicate the state to all nodes
func (s *server) broadcastState() {
	// Lock on the serves is assumed to be held by the caller
	if s.leaderId != s.nodeId {
		// Only the leader can broadcast the state
		return
	}

	log.Printf("Node %d: Broadcasting state to all nodes with lower id.", s.nodeId)
	for n := 0; n < s.nodeId; n++ {
		if n == s.nodeId {
			continue
		}

		node := fmt.Sprintf("localhost:%d", *basePort+n)

		conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
		if err != nil {
			log.Printf("Node %d: Could not connect to node %s", s.nodeId, node)
			continue
		}
		defer conn.Close()

		c := pb.NewAuctionClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err = c.SyncState(ctx, &pb.SyncStateRequest{HighestBid: s.highestBid, HighestBidderId: s.highestBidderId})
		if err != nil {
			log.Printf("Node %d: Could not sync leader to node %s", s.nodeId, node)
			continue
		}

		log.Printf("Node %d: Synced state to node %s", s.nodeId, node)
	}

	log.Printf("Node %d: Synced state to all active nodes", s.nodeId)
}

// Remote called by a node with a lower id as a part of the election process
func (s *server) SendElection(ctx context.Context, req *pb.ElectionRequest) (*pb.ElectionResponse, error) {
	log.Printf("Node %d: SendElection: %d", s.nodeId, req.NodeId)
	if req.NodeId < int32(s.nodeId) {
		s.startElection()
		return &pb.ElectionResponse{Ok: true}, nil
	}

	return &pb.ElectionResponse{Ok: false}, nil
}

// Remote called by the leader node to announce the new leader to all nodes with lower id
func (s *server) AnnounceLeader(ctx context.Context, req *pb.AnnounceLeaderRequest) (*pb.Empty, error) {
	log.Printf("Node %d: Updating leader: %d", s.nodeId, req.LeaderId)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.leaderId = int(req.LeaderId)
	log.Printf("Node %d: New leader is %d", s.nodeId, s.leaderId)
	return &pb.Empty{}, nil
}

// Remote called by the leader node to sync the state to all nodes with lower id
// Called when the leader changes and needs to replicate the state to all nodes
func (s *server) SyncState(ctx context.Context, req *pb.SyncStateRequest) (*pb.Empty, error) {
	log.Printf("Node %d: Updating state HighestBid: %d, HighestBidderId: %d", s.nodeId, req.HighestBid, req.HighestBidderId)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.highestBid = req.HighestBid
	s.highestBidderId = req.HighestBidderId
	log.Printf("Node %d: Updated state HighestBid: %d, HighestBidderId: %d", s.nodeId, req.HighestBid, req.HighestBidderId)

	return &pb.Empty{}, nil
}

func runAuction(s *server) {
	s.auctionState = pb.AUCTION_STATUS_ONGOING
	// Sleep for the duration of the auction
	time.Sleep(time.Minute * time.Duration(*min))
	s.auctionState = pb.AUCTION_STATUS_FINISHED
}

func main() {
	fmt.Println("Starting node...")

	flag.Parse()
	if *nodeId < 0 {
		log.Fatalf("nodeId must be specified")
	}

	if *nodeId >= *maxNode {
		log.Fatalf("nodeId must be lower than maxNodeId")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *basePort+*nodeId))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	var server *server = &server{
		auctionState:    pb.AUCTION_STATUS_NOTSTARTED,
		highestBid:      0,
		highestBidderId: -1,
		nodeId:          *nodeId,
		leaderId:        -1,
	}

	go runAuction(server)

	pb.RegisterAuctionServer(s, server)

	log.Printf("Node listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
