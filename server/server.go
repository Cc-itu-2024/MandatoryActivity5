package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"net"
	"sync"

	pb "auction/auction"

	"google.golang.org/grpc"
)

var server *Server
var ports = []string(":5050", ":5051", ":5052", ":5053")
var Auctionactive bool

type AuctionServer struct {
	pb.UnimplementedAuctionServer
	client pb.AuctionClient
	serverNodes map[int32]
	mu sync.Mutex
	isLeader bool	
}

type Node struct {

}

func NewAuctionServer() * AuctionServer{

}

func main() {
	logFile, err := os.OpenFile("auction.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server = new Server(ports[0])

	grpcServer := grpc.NewServer()
	pb.RegisterAuctionServer(grpcServer, server)

	log.Printf("Server listening on %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


func runAuction(){



}

func startServer(){


}


func bid(){

}

func ack(){

}

func outcome(){

}


