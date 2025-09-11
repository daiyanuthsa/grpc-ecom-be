package jwt

import (
	"context"
	"fmt"
	"os"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	RoleCode string `json:"role_code"`
}
type JwtEntityContextKey string
var JwtEntityContextValue JwtEntityContextKey = "JwtEntity"

func (jc *JWTClaims) SetToContext(ctx context.Context) context.Context{
	ctx = context.WithValue(ctx, JwtEntityContextValue, jc)
	return ctx
}

func GetClaimsFromToken(token string) (*JWTClaims, error) {

	tokenClaims, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected token signing method %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}
	if !tokenClaims.Valid {
		return nil, utils.UnauthenticatedResponse()
	}

	claims, ok := tokenClaims.Claims.(*JWTClaims)
	if ok {
		return claims, nil
	}
	return nil, utils.UnauthenticatedResponse()
}