package service

import (
	"context"
	"time"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/entity"
	jwtentity "github.com/daiyanuthsa/grpc-ecom-be/internal/entity/jwt"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/repository"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/utils"
	"github.com/daiyanuthsa/grpc-ecom-be/pb/product"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IProductService interface {
	CreateProduct(ctx context.Context, request *product.CreateProductRequest) (*product.CreateProductResponse, error)
}

type productService struct {
	productRepository repository.IProductRepository
	storageService IStorageService
}

func (ps *productService) CreateProduct(ctx context.Context, request *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	// cek apakah benar admin
	print("cek admin")
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}	

	if claims.RoleCode != entity.UserRoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}
	// TODO: Cek apakah request.ImageFileName tersedia
	objectKey := request.ImageFileName
	exists, err := ps.storageService.CheckIfObjectExists(ctx, objectKey)
    if err != nil {
        // Error koneksi atau sistem R2
        return nil, status.Errorf(codes.Internal, "Storage service check failed: %v", err)
    }
	if !exists {
        // File tidak ditemukan di R2 (Ini adalah error bisnis/validasi)
        return &product.CreateProductResponse{
            Base: utils.BadRequestResponse("Image file not found in storage. Please upload the image first."),
        }, nil
    }

	newProduct := entity.Product{
		Id:            uuid.NewString(),
		Name:          request.Name,
		Description:   request.Description,
		Price:         float64(request.Price),
		ImageFileName: request.ImageFileName,
		CreatedAt:     time.Now(),
		CreatedBy:     &claims.FullName,
	}
	productId := newProduct.Id
	err = ps.productRepository.CreateProduct(ctx, &newProduct)
	if err != nil {
		return nil, err
	}

	return &product.CreateProductResponse{
		Base: utils.SuccessResponse("Product created successfully"),
		Id: productId,
	}, nil
}

func NewProductService(productRepository repository.IProductRepository, storageService IStorageService) IProductService {
	return &productService{
		productRepository: productRepository,
		storageService: storageService,
	}
}