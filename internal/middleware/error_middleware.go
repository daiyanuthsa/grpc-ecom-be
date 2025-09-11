package middleware

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler)(res any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Internal server error")
		}
	}()

	// Panggil handler dulu
	res, err = handler(ctx, req)
	
	// Check if an error was returned.
	if err != nil {
		// Log the original error for debugging purposes.
		log.Printf("gRPC error: %v", err)

		// Check if the error is already a gRPC status error.
		// If it is (like Unauthenticated), we don't need to do anything.
		// We just let it pass through to the client.
		if _, ok := status.FromError(err); ok {
			return nil, err
		}

		// If it's NOT a gRPC status error (e.g., a database error like "sql: no rows"),
		// then we wrap it in a generic Internal error so we don't leak implementation details to the client.
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	// If there was no error, return the response as normal.
	return res, nil
}