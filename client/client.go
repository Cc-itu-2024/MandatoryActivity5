package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "auction/auction"

	"google.golang.org/grpc"
)

const (
	timeoutDuration = 5 * time.Second // Timeout for gRPC requests
)

// Connect to the auction server
func connectToServer(address string) (pb.AuctionClient, *grpc.ClientConn, error) {
	// Set up a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(timeoutDuration))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	// Create a new Auction client
	client := pb.NewAuctionClient(conn)
	return client, conn, nil
}

// Send a bid to the auction server
func bid(client pb.AuctionClient, amount int32) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	response, err := client.Bid(ctx, &pb.Amount{Amount: amount})
	if err != nil {
		log.Printf("Error while sending bid: %v", err)
		return
	}

	fmt.Printf("Bid response: %s\n", response.Message)
}

// Get the current auction result
func getAuctionResult(client pb.AuctionClient) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	response, err := client.Result(ctx, &pb.Void{})
	if err != nil {
		log.Printf("Error while querying result: %v", err)
		return
	}

	if response.HighestBidder == "in progress" {
		fmt.Printf("Auction1 is still in progress. Highest bid so far: %d\n", response.HighestBid)
	} else {
		fmt.Printf("Auction1 ended. Winner: %s with bid: %d\n", response.HighestBidder, response.HighestBid)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client.go <server_address>")
		return
	}

	// Server address from command line argument
	serverAddress := os.Args[1]

	// Connect to the auction server
	client, conn, err := connectToServer(serverAddress)
	if err != nil {
		log.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connected to auction server.")

	for {
		fmt.Println("\n1. Place a Bid")
		fmt.Println("2. Get Auction Result")
		fmt.Println("3. Exit")

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			var bidAmount int32
			fmt.Print("Enter your bid amount: ")
			fmt.Scanln(&bidAmount)
			bid(client, bidAmount)

		case 2:
			getAuctionResult(client)

		case 3:
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
