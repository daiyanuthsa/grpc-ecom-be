package handler

import (
	"context"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/service"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/utils"

	"github.com/daiyanuthsa/grpc-ecom-be/pb/product"
)



type productHandler struct {
	product.UnimplementedProductServiceServer
	
	productService service.IProductService
}

func (ph *productHandler) CreateProduct(ctx context.Context, request *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return &product.CreateProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ph.productService.CreateProduct(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewProductHandler(productService service.IProductService) *productHandler {
	return &productHandler{
		productService: productService,
	}
}
