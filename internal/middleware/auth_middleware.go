package middleware

import (
	"context"
	"log"

	jwtentity "github.com/daiyanuthsa/grpc-ecom-be/internal/entity/jwt"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/utils"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)


type authMiddleware struct{
	cacheService *gocache.Cache
	whitelist    map[string]struct{} // The set of whitelisted endpoints
}
func (am *authMiddleware) Middleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler)(res any, err error) {
	
	log.Println(info.FullMethod)
	// 💡 Check if the method is in our whitelist map.
    if _, ok := am.whitelist[info.FullMethod]; ok {
        // If it is, skip all auth checks and proceed directly to the handler.
        return handler(ctx, req)
    }
	// Ambil token dari meta data
	tokenStr, err :=jwtentity.ParseTokenFromContex(ctx)
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}

	// Cek blacklist/logout token dari cache
	 _, ok := am.cacheService.Get(tokenStr)
	 if ok{
			// Kalau ketemu di cache berarti token sudah di logout
			return nil, utils.UnauthenticatedResponse()
		}

	//parse token menjadi entity.jwt
	claims, err := jwtentity.GetClaimsFromToken(tokenStr)
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}

	//sematkan entity.jwt ke contex
	ctx = claims.SetToContext(ctx)

	// Panggil handler dulu
	res, err = handler(ctx, req)
	
	return res, err
}

func NewAuthMiddleware(cacheService *gocache.Cache,  publicEndpoints []string)*authMiddleware{
	// Create the whitelist map for efficient lookups
    whitelist := make(map[string]struct{})
    for _, endpoint := range publicEndpoints {
        whitelist[endpoint] = struct{}{}
    }

	return &authMiddleware {
		cacheService: cacheService,
		whitelist:    whitelist,
	}
}