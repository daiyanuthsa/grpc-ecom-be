package handler

import (
	"context"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/service"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/utils"

	"github.com/daiyanuthsa/grpc-ecom-be/pb/auth"
)



type authHandler struct {
	auth.UnimplementedAuthServiceServer

	authService service.IAuthService
}

func (ah *authHandler) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {

		return &auth.RegisterResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ah.authService.Register(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ah *authHandler) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {

		return &auth.LoginResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ah.authService.Login(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ah *authHandler) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	// Implement your logout logic here (if any)
	res, err := ah.authService.Logout(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ah *authHandler) ChangePassword(ctx context.Context, request *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return &auth.ChangePasswordResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ah.authService.ChangePassword(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (ah *authHandler) GetProfile(ctx context.Context, request *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	
	res, err := ah.authService.GetProfile(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func NewAuthHandler(authService service.IAuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}
