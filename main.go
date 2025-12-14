package main

import (
	"log"
	"net"
	"os"

	checker "github.com/impruthvi/pulse-check-apis/checker/v1"
	monitor "github.com/impruthvi/pulse-check-apis/monitor/v1"
	"github.com/impruthvi/pulse-check-monitor/db"
	"github.com/impruthvi/pulse-check-monitor/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB_URL is required")
	}

	dbProvider := db.New(dbURL)

	server := grpc.NewServer()

	checkerServiceURL := os.Getenv("CHECKER_SERVICE_URL")
	if checkerServiceURL == "" {
		log.Fatal("CHECKER_SERVICE_URL is required")
	}

	grpcConn, err := grpc.NewClient(checkerServiceURL, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))

	if err != nil {
		log.Fatalf("Failed to dial CheckService: %v", err)
	}

	defer grpcConn.Close()

	checkerServiceClient := checker.NewCheckerServiceClient(grpcConn)

	monitorService := service.New(
		service.Dependencies{
			DBProvider:    dbProvider,
			CheckerClient: checkerServiceClient,
		},
	)

	monitor.RegisterMonitorServiceServer(server, monitorService)

	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC Server is running on port 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
