package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/handler"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/repository"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/service"
	"github.com/daiyanuthsa/grpc-ecom-be/pb/auth"
	"github.com/daiyanuthsa/grpc-ecom-be/pb/product"
	gocache "github.com/patrickmn/go-cache"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/middleware"
	database "github.com/daiyanuthsa/grpc-ecom-be/pkg"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	godotenv.Load()

	log.Println("Starting gRPC server...")

	// Create a listener on TCP port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := database.ConnectDB(ctx, os.Getenv("DB_URI"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	// Define all public (unauthenticated) endpoints in a single, clean slice.
    publicEndpoints := []string{
        "/auth.AuthService/Login",
        "/auth.AuthService/Register", // Easy to add new ones!
		"/product.ProductService/DetailProduct",
    }

	cacheService := gocache.New(time.Hour*24, time.Hour)
	authMiddleware := middleware.NewAuthMiddleware(cacheService, publicEndpoints)
	authService := service.NewAuthService(repository.NewAuthRepository(db), cacheService)
	authHandler := handler.NewAuthHandler(authService)

	productService := service.NewProductService(repository.NewProductRepository(db), service.NewStorageService(ctx))
	productHandler := handler.NewProductHandler(productService)
	

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.ErrorMiddleware,
			authMiddleware.Middleware), // Add the error handling middleware
	)

	auth.RegisterAuthServiceServer(serv, authHandler)
	product.RegisterProductServiceServer(serv, productHandler)

	if os.Getenv("ENVIRONMENT") == "dev" {
		reflection.Register(serv)
		log.Println("Reflection service registered")
	}

	log.Println("Server is running on port 50051")
	if err := serv.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}

}