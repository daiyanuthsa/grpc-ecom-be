package service

import (
	"context"
	"fmt"
	"os"
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
    DetailProduct(ctx context.Context, request *product.DetailProductRequest) (*product.DetailProductResponse, error)
	UpdateProduct(ctx context.Context, request *product.UpdateProductRequest) (*product.UpdateProductResponse, error)
}

type productService struct {
	productRepository repository.IProductRepository
	storageService IStorageService
}

func (ps *productService) CreateProduct(ctx context.Context, request *product.CreateProductRequest) (*product.CreateProductResponse, error) {
	// cek apakah benar admin
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

func (ps *productService) DetailProduct(ctx context.Context, request *product.DetailProductRequest) (*product.DetailProductResponse, error){
	productData, err := ps.productRepository.GetProductById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if productData == nil {
		return &product.DetailProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	return &product.DetailProductResponse{
		Base: utils.SuccessResponse("Product retrieved successfully"),
		Id:            productData.Id,
		Name:          productData.Name,
		Description:   productData.Description,
		Price:         productData.Price,
		ImageUrl: 		fmt.Sprintf("%s/%s", os.Getenv("R2_PUBLIC_DOMAIN"), productData.ImageFileName),
	}, nil
}

func (ps *productService) UpdateProduct(ctx context.Context, request *product.UpdateProductRequest) (*product.UpdateProductResponse, error){
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}	

	if claims.RoleCode != entity.UserRoleAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied")
	}
	// Check is the product exist
	productData, err := ps.productRepository.GetProductById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if productData == nil {
		return &product.UpdateProductResponse{
			Base: utils.NotFoundResponse("Product not found"),
		}, nil
	}

	if request.ImageFileName != productData.ImageFileName {
		objectKey := request.ImageFileName
		exists, err := ps.storageService.CheckIfObjectExists(ctx, objectKey)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Storage service check failed: %v", err)
		}
		if !exists {
			return &product.UpdateProductResponse{
				Base: utils.NotFoundResponse("New image file not found in storage. Please upload the image first."),
			}, nil
		}
		// If the image file name has changed, delete the old image from storage
		if productData.ImageFileName != ""  {
			delErr := ps.storageService.DeleteObject(ctx, productData.ImageFileName)
			if delErr != nil {
				// Log the error but don't block the product update if image deletion fails
				fmt.Printf("Failed to delete old image %s from storage: %v\n", productData.ImageFileName, delErr)
			}
		}
		
	}

	productData.Name = request.Name
	productData.Description = request.Description
	productData.Price = float64(request.Price)
	productData.ImageFileName = request.ImageFileName
	productData.UpdatedAt = time.Now()
	productData.UpdatedBy = &claims.FullName

	err = ps.productRepository.UpdateProduct(ctx, productData)
	if err != nil {
		return nil, err
	}

	return &product.UpdateProductResponse{
		Base: utils.SuccessResponse("Product updated successfully"),
		Id: productData.Id,
	}, nil
}

func NewProductService(productRepository repository.IProductRepository, storageService IStorageService) IProductService {
	return &productService{
		productRepository: productRepository,
		storageService: storageService,
	}
}