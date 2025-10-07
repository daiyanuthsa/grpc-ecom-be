package service

import (
	"context"
	"log"
	"time"

	"github.com/daiyanuthsa/grpc-ecom-be/internal/entity"
	jwtentity "github.com/daiyanuthsa/grpc-ecom-be/internal/entity/jwt"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/repository"
	"github.com/daiyanuthsa/grpc-ecom-be/internal/utils"
	"github.com/daiyanuthsa/grpc-ecom-be/pb/cart"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ICartService defines the interface for cart-related business logic.
type ICartService interface {
	// AddProductToCart adds a product to the user's cart.
	AddProductToCart(ctx context.Context, req *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error)
	// ListCart retrieves all items in the user's cart.
	ListCart(ctx context.Context, request *cart.ListCartRequest) (*cart.ListCartResponse, error)
	// UpdateCartItem updates the quantity of a specific cart item.
	UpdateCartItem(ctx context.Context, request *cart.UpdateCartItemRequest) (*cart.UpdateCartItemResponse, error)
	// DeleteCartItem removes a specific cart item from the cart.
	DeleteCartItem(ctx context.Context, request *cart.DeleteCartItemRequest) (*cart.DeleteCartItemResponse, error)
}

// CartService implements ICartService.
type CartService struct {
	cartRepository    repository.ICartRepository
	productRepository repository.IProductRepository
}

// NewCartService creates a new instance of CartService.
func NewCartService(cartRepository repository.ICartRepository, productRepository repository.IProductRepository) ICartService {
	return &CartService{
		cartRepository:    cartRepository,
		productRepository: productRepository,
	}
}

// AddProductToCart adds a product to the user's cart or updates its quantity if already present.
func (s *CartService) AddProductToCart(ctx context.Context, req *cart.AddProductToCartRequest) (*cart.AddProductToCartResponse, error) {
	// cek apakah id produk ada di DB atau tidak (gunakan produk repository)
	product, err := s.productRepository.GetProductById(ctx, req.ProductId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if product == nil {
		return &cart.AddProductToCartResponse{
			Base: utils.NotFoundResponse("product not found"),
		}, nil
	}

	// cek siapa usernya
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}	

	// cek apakah procuk sudah ada di cart user
	cartItem, err := s.cartRepository.FindByUserIDAndProductID(ctx, claims.Subject, req.ProductId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var cartID uuid.UUID

	if cartItem != nil {
		//  sudah ada -> update
		cartItem.Quantity += 1 // Cast to int64
		now := time.Now()
		updatedBy := claims.FullName
		cartItem.UpdatedAt = &now
		cartItem.UpdatedBy = &updatedBy

		err = s.cartRepository.Update(ctx, cartItem)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		cartID = cartItem.ID
	} else {
		//belum ada -> insert
		newItem := &entity.CartItem{
			ID:        uuid.New(), // Generate new UUID
			UserID:    claims.Subject,
			ProductID: req.ProductId,
			Quantity:  1, // Cast to int64
			CreatedAt: time.Now(),
			CreatedBy: claims.FullName,
		}
		err = s.cartRepository.Insert(ctx, newItem)
		if err != nil {
			return &cart.AddProductToCartResponse{
				Base: utils.BadRequestResponse("Failed to add product to cart"),
				Id:   "",
			}, nil
		}
		cartID = newItem.ID
	}

	return &cart.AddProductToCartResponse{
		Base: utils.SuccessResponse("Product added to cart successfully"),
		Id:   cartID.String(),
	}, nil
}

// ListCart retrieves all items in the user's cart along with product details and total price.
func (s *CartService) ListCart(ctx context.Context, request *cart.ListCartRequest) (*cart.ListCartResponse, error){
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}

	cartItems, err := s.cartRepository.FindByUserID(ctx, claims.Subject)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var responseItems []*cart.CartItem
	var totalPrice float64

	for _, item := range cartItems {
		product, err := s.productRepository.GetProductById(ctx, item.ProductID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get product details for product ID %s: %v", item.ProductID, err)
		}
		if product == nil {
			log.Printf("WARNING: Product with ID %s found in cart but not in products collection", item.ProductID)
			continue
		}

		responseItems = append(responseItems, &cart.CartItem{
			CartId:      item.ID.String(),
			ProductId:   product.Id,
			ProductName: product.Name,
			ImageUrl:    product.ImageFileName,
			Price:       product.Price,
			Quantity:    int32(item.Quantity),
		})
		totalPrice += product.Price * float64(item.Quantity)
	}

	return &cart.ListCartResponse{
		Base:      utils.SuccessResponse("Cart items retrieved successfully"),
		Items:     responseItems,
		TotalPrice: totalPrice,
	}, nil
}

// UpdateCartItem updates the quantity of a specific cart item for the authenticated user.
func (s *CartService) UpdateCartItem(ctx context.Context, request *cart.UpdateCartItemRequest) (*cart.UpdateCartItemResponse, error){
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}

	cartUUID, err := uuid.Parse(request.CartId)
	if err != nil {
		return &cart.UpdateCartItemResponse{
			Base: utils.BadRequestResponse("Invalid cart ID format"),
		}, nil
	}

	cartItem, err := s.cartRepository.FindByID(ctx, cartUUID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if cartItem == nil {
		return &cart.UpdateCartItemResponse{
			Base: utils.NotFoundResponse("Cart item not found"),
		}, nil
	}

	if cartItem.UserID != claims.Subject {
		return nil, utils.UnauthenticatedResponse()
	}

	cartItem.Quantity = int(request.NewQuantity) 
	now := time.Now()
	updatedBy := claims.FullName
	cartItem.UpdatedAt = &now
	cartItem.UpdatedBy = &updatedBy

	err = s.cartRepository.Update(ctx, cartItem)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cart.UpdateCartItemResponse{
		Base: utils.SuccessResponse("Cart item updated successfully"),
	}, nil
}

// DeleteCartItem removes a specific cart item from the authenticated user's cart.
func (s *CartService) DeleteCartItem(ctx context.Context, request *cart.DeleteCartItemRequest) (*cart.DeleteCartItemResponse, error){
	claims, err := jwtentity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, utils.UnauthenticatedResponse()
	}

	cartUUID, err := uuid.Parse(request.CartId)
	if err != nil {
		return &cart.DeleteCartItemResponse{
			Base: utils.BadRequestResponse("Invalid cart ID format"),
		}, nil
	}

	cartItem, err := s.cartRepository.FindByID(ctx, cartUUID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if cartItem == nil {
		return &cart.DeleteCartItemResponse{
			Base: utils.NotFoundResponse("Cart item not found"),
		}, nil
	}

	if cartItem.UserID != claims.Subject {
		return nil, utils.UnauthenticatedResponse()													
	}

	err = s.cartRepository.Delete(ctx, cartUUID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &cart.DeleteCartItemResponse{
		Base: utils.SuccessResponse("Cart item deleted successfully"),
	}, nil
}
