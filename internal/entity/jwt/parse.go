package jwt

import (
	"context"
	"strings"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/utils"
	"google.golang.org/grpc/metadata"
)

func ParseTokenFromContex(ctx context.Context)(string,error){
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", utils.UnauthenticatedResponse()
	}
	bearerToken := md["authorization"]
	if len(bearerToken) == 0 {
		return "", utils.UnauthenticatedResponse()
	}
	
	tokenSSplit := strings.SplitN(bearerToken[0], " ", 2)

	if len(tokenSSplit) != 2 || tokenSSplit[0] != "Bearer" {
		return "", utils.UnauthenticatedResponse()
	}

	return tokenSSplit[1], nil
}