package handler

import (
	"context"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/service"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/utils"

	"github.com/daiyanuthsa/grpc-ecom-be/pb/cart"
)



type cartHandler struct {
	cart.UnimplementedCartServiceServer

	cartService service.ICartService
}

func (ch *cartHandler) AddProductToCart(ctx context.Context, request *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return &cart.AddProductToCartResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ch.cartService.AddProductToCart(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ch *cartHandler) ListCart(ctx context.Context, request *cart.ListCartRequest) (*cart.ListCartResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return &cart.ListCartResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ch.cartService.ListCart(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ch *cartHandler) UpdateCartItem(ctx context.Context, request *cart.UpdateCartItemRequest) (*cart.UpdateCartItemResponse, error){
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return &cart.UpdateCartItemResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ch.cartService.UpdateCartItem(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ch *cartHandler) DeleteCartItem(ctx context.Context, request *cart.DeleteCartItemRequest) (*cart.DeleteCartItemResponse, error){
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return &cart.DeleteCartItemResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ch.cartService.DeleteCartItem(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}


func NewCartHandler(cartService service.ICartService) *cartHandler {
	return &cartHandler{
		cartService: cartService,
	}
}
