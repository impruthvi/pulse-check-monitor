package main

import (
	"context"
	"log"
	"net"
	"os"

	checker "github.com/impruthvi/pulse-check-apis/checker/v1"
	monitor "github.com/impruthvi/pulse-check-apis/monitor/v1"
	"github.com/impruthvi/pulse-check-monitor/db"
	"github.com/impruthvi/pulse-check-monitor/service"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
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

	ctx := context.Background()

	otlpEndpoint := os.Getenv("OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "localhost:4317"
	}

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otlpEndpoint),
		otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		log.Printf("Warning: Failed to create OTLP trace exporter: %v", err)
		log.Println("Continuing without tracing...")
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("monitord"),
		)),
	)

	otel.SetTracerProvider(tracerProvider)
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	checkerServiceURL := os.Getenv("CHECKER_SERVICE_URL")
	if checkerServiceURL == "" {
		log.Fatal("CHECKER_SERVICE_URL is required")
	}

	openTelemetryClientHandler := otelgrpc.NewClientHandler(
		otelgrpc.WithTracerProvider(tracerProvider),
	)

	grpcConn, err := grpc.NewClient(
		checkerServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(openTelemetryClientHandler),
	)

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
