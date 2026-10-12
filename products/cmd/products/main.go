package main

import (
	"log"
	"net"

	pb "hw2/products/gen/go/proto"
	"hw2/products/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	server := grpc.NewServer()
	reflection.Register(server)
	
	service, err := service.NewProductService()
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
		return
	}
	
	pb.RegisterProductServiceServer(server, service)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen port: %v", err)
	}

	log.Println("Listen RPC on :50051")
	log.Fatal(server.Serve(lis))
}