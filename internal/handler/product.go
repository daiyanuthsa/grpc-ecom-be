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

func (ph *productHandler) DetailProduct(ctx context.Context, request *product.DetailProductRequest) (*product.DetailProductResponse, error){
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return &product.DetailProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ph.productService.DetailProduct(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (ph *productHandler) UpdateProduct(ctx context.Context, request *product.UpdateProductRequest) (*product.UpdateProductResponse, error){
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return &product.UpdateProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ph.productService.UpdateProduct(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ph *productHandler) DeleteProduct(ctx context.Context, request *product.DeleteProductRequest) (*product.DeleteProductResponse, error){
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return &product.DeleteProductResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ph.productService.DeleteProduct(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ph *productHandler) ListProducts(ctx context.Context, request *product.ListProductsRequest) (*product.ListProductsResponse, error){
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return &product.ListProductsResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ph.productService.ListProducts(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ph *productHandler) ListProductsAdmin(ctx context.Context, request *product.ListProductsAdminRequest) (*product.ListProductsAdminResponse, error){
validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if len(validationErrors) > 0 {
		return &product.ListProductsAdminResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	res, err := ph.productService.ListProductsAdmin(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ph *productHandler) HighlightProducts(ctx context.Context, request *product.HighlightProductsRequest) (*product.HighlightProductsResponse, error){
	res, err := ph.productService.HighlightProducts(ctx, request)
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
