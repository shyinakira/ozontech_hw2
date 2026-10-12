package main

import (
	"log"
	"net"
	"net/http"

	"hw2/cart/internal/handlers"
	"hw2/cart/internal/storage"
	pb "hw2/products/gen/go/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen :8080, %v", err)
	}

	
	storage := storage.NewStorage()

	connGrpc, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to create grpc conn :50051, %v\n", err)
		// log.Fatalf("failed to create grpc conn :50051, %v", err)
	}
	productsClient := pb.NewProductServiceClient(connGrpc)

	cartHandler := handlers.NewCartHandler(storage, productsClient)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /user/{user_id}/cart/{sku}", cartHandler.AddCartItem)
	mux.HandleFunc("GET /user/{user_id}/cart", cartHandler.GetCart)

	log.Println("start :8080 http server")
	log.Fatalln(http.Serve(lis, mux))
}