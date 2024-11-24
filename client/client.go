package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	pb "auction/auction"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	nodeId    = flag.Int("nodeId", 0, "The nodeId")
	//maxNode = flag.Int("maxNode", -1, "The maximum number of nodes")
	basePort  = flag.Int("basePort", 50000, "The base port number. The server will listen on this port + nodeId")
	bidderId  = flag.Int("bidderId", 0, "The bidder id")
)

func main() {

	flag.Parse()

	//if *maxNode == -1 {
	//	log.Fatalf("maxNodeId is required")
	//}

	// Set up a connection to the server.
	addr := fmt.Sprintf("localhost:%d", *basePort+*nodeId)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewAuctionClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reader := bufio.NewReader(os.Stdin)
	for {
		log.Println("Place bid:")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg) // Get rid of the newline character

		if msg == "/exit" {
			log.Println("Exiting...")
			break
		}

		bidValue, err := strconv.Atoi(msg)
		if err != nil {
			log.Fatalf("Invalid bid value: %v", err)
		}

		response, err := client.Bid(ctx, &pb.BidRequest{Bid: int32(bidValue), BidderId: int32(*bidderId)})
		if err != nil {
			log.Fatalf("Bid(%v) failed: %v", msg, err)
		}
		log.Println("Response: ", response.Status)

		result, err := client.Result(ctx, &pb.Empty{})
		if err != nil {
			log.Fatalf("Result() failed: %v", err)
		}

		log.Printf("Highest bid: %d by bidder %d", result.HighestBid, result.BidderId)
	}
}
